package http

import (
	"net/http"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/ports"
	httpSwagger "github.com/swaggo/http-swagger"
)

// NewRouter creates and configures the HTTP router with all routes and middleware
func NewRouter(orderService ports.OrderService) http.Handler {
	mux := http.NewServeMux()

	// Initialize handlers
	orderHandler := NewOrderHandler(orderService)

	// Health check
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		RespondJSON(w, http.StatusOK, map[string]string{
			"status": "ok",
		})
	})

	// Swagger documentation
	mux.HandleFunc("GET /swagger/", httpSwagger.WrapHandler)

	// API v1 routes - Orders
	mux.HandleFunc("POST /api/v1/orders", orderHandler.CreateOrder)
	mux.HandleFunc("GET /api/v1/orders/{code}/total", orderHandler.GetOrderTotal)

	// API v1 routes - Customers
	mux.HandleFunc("GET /api/v1/customers/{code}/orders", orderHandler.ListCustomerOrders)
	mux.HandleFunc("GET /api/v1/customers/{code}/orders/count", orderHandler.CountCustomerOrders)

	// Apply middleware chain
	return Chain(mux,
		LoggingMiddleware,
		CORSMiddleware,
	)
}
