package http

import (
	"log"
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Initialize handlers
	orderHandler := NewOrderHandler()

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

	return logMiddleware(mux)
}

func logMiddleware(next *http.ServeMux) *http.ServeMux {
	wrapper := http.NewServeMux()
	
	wrapper.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.Method, r.URL.Path, r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
	
	return wrapper
}
