package server

import (
	"fmt"
	"net/http"
	"time"

	"go.uber.org/zap"

	httphandler "github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/inbound/http"
	db "github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/database"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/application/services"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/config"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/logger"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/ports"
)

type Server struct {
	router       http.Handler
	server       *http.Server
	orderService ports.OrderService
}

func NewServer(cfg *config.Config, dbStore *db.Store, messagePublisher ports.MessagePublisher) *Server {
	// Initialize service with dependency injection
	orderService := services.NewOrderService(dbStore, messagePublisher)

	// Initialize router with service
	router := httphandler.NewRouter(cfg, orderService)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		router:       router,
		server:       server,
		orderService: orderService,
	}
}

func (s *Server) Start() error {
	logger.Info("HTTP server starting",
		zap.String("address", s.server.Addr),
		zap.Float64("read_timeout", s.server.ReadTimeout.Seconds()),
		zap.Float64("write_timeout", s.server.WriteTimeout.Seconds()),
	)

	logger.Info("Available endpoints",
		zap.String("health", "GET /health"),
		zap.String("metrics", "GET /metrics"),
		zap.String("swagger", "GET /swagger/index.html"),
		zap.String("order_total", "GET /api/v1/orders/{code}/total"),
		zap.String("customer_orders", "GET /api/v1/customers/{code}/orders"),
		zap.String("customer_orders_count", "GET /api/v1/customers/{code}/orders/count"),
		zap.String("create_order", "POST /api/v1/orders"),
	)

	logger.Info("OrderService initialized", zap.String("status", "ready"))

	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	logger.Info("Server shutdown initiated")
	return s.server.Close()
}
