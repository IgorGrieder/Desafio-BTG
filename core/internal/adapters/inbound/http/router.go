package http

import (
	"net/http"
	"time"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/inbound/http/middleware"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/config"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/infrastructure/telemetry"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/ports"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/pkg/httputils"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// spanNames maps route patterns to custom span names
var spanNames = map[string]string{
	"GET /health":                               "health",
	"GET /metrics":                              "metrics",
	"POST /api/v1/orders":                       "orders.create",
	"GET /api/v1/orders/{code}/total":           "orders.getTotal",
	"GET /api/v1/customers/{code}/orders":       "customers.listOrders",
	"GET /api/v1/customers/{code}/orders/count": "customers.countOrders",
}

// NewRouter creates and configures the HTTP router with all routes and middleware
func NewRouter(cfg *config.Config, orderService ports.OrderService) http.Handler {
	mux := http.NewServeMux()

	// Initialize handlers
	orderHandler := NewOrderHandler(orderService)
	healthHandler := NewHealthHandler()

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		httputils.RespondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
			"time":   time.Now().Format(time.RFC3339),
		})
	})

	// Metrics endpoint
	mux.Handle("GET /metrics", healthHandler.Metrics())

	// Swagger documentation
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// API v1 routes - Orders
	mux.HandleFunc("POST /api/v1/orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /api/v1/orders/{code}/total", orderHandler.GetOrderTotal)

	// API v1 routes - Customers
	mux.HandleFunc("GET /api/v1/customers/{code}/orders", orderHandler.ListCustomerOrders)
	mux.HandleFunc("GET /api/v1/customers/{code}/orders/count", orderHandler.CountCustomerOrders)

	// Wrap with global middlewares: metrics -> logging -> CORS -> routes
	innerHandler := middleware.MetricsMiddleware(
		middleware.LoggingMiddleware(
			middleware.CORSMiddleware(mux),
		),
	)

	// Wrap with otelhttp for automatic tracing with custom span names
	otelOptions := []otelhttp.Option{
		otelhttp.WithSpanNameFormatter(func(operation string, r *http.Request) string {
			key := r.Method + " " + r.Pattern
			if name, ok := spanNames[key]; ok {
				return name
			}
			if r.Pattern != "" {
				return r.Pattern
			}
			return r.URL.Path
		}),
	}

	// Only add tracer provider if telemetry is initialized
	if telemetry.TracerProvider != nil {
		otelOptions = append(otelOptions, otelhttp.WithTracerProvider(telemetry.TracerProvider))
	}

	handler := otelhttp.NewHandler(
		innerHandler,
		cfg.App.Name,
		otelOptions...,
	)

	return handler
}
