package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	httpAdapter "github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/inbound/http"
	db "github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/database"
	database "github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/database/sqlc"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/messaging"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/application/services"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/config"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/infrastructure/telemetry"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/logger"

	_ "github.com/IgorGrieder/Desafio-BTG/tree/main/core/docs"
)

// @title BTG Pactual Order Processing API
// @version 1.0
// @description API for managing and querying orders
// @host localhost:8080
// @BasePath /api/v1
func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load configuration: %v\n", err)
		os.Exit(1)
	}

	// Initialize logger
	if err := logger.Init(cfg.App.Env); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("Starting application",
		zap.String("name", cfg.App.Name),
		zap.String("version", cfg.App.Version),
		zap.String("env", cfg.App.Env),
	)

	// Initialize OpenTelemetry tracer
	var shutdownTracer func(context.Context) error
	if cfg.OTel.Enabled {
		var err error
		shutdownTracer, err = telemetry.InitTracer(
			cfg.OTel.Endpoint,
			cfg.App.Name,
			cfg.App.Version,
		)
		if err != nil {
			logger.Warn("Failed to initialize tracer, continuing without tracing", zap.Error(err))
		} else {
			logger.Info("OpenTelemetry tracer initialized", zap.String("endpoint", cfg.OTel.Endpoint))
		}
	}

	// Initialize database connection
	dbConn, err := db.NewDB(ctx, cfg.Database.DSN())
	if err != nil {
		logger.Fatal("Failed to connect to database",
			zap.Error(err),
		)
	}
	defer dbConn.Close()

	logger.Info("Connected to PostgreSQL")

	// Initialize SQLC queries
	queries := database.New(dbConn.Pool)
	dbStore := db.NewStore(dbConn, queries)

	logger.Info("Repository established")

	// Initialize RabbitMQ publisher
	publisher, err := messaging.NewRabbitMQPublisher(
		cfg.RabbitMQ.URL(),
		cfg.RabbitMQ.Exchange,
		cfg.RabbitMQ.Queue,
	)
	if err != nil {
		logger.Fatal("Failed to initialize RabbitMQ publisher",
			zap.Error(err),
		)
	}
	defer publisher.Close()

	logger.Info("RabbitMQ publisher initialized")

	// Initialize services
	orderService := services.NewOrderService(dbStore, publisher)

	// Initialize HTTP router with middleware chain
	router := httpAdapter.NewRouter(cfg, orderService)

	// Create HTTP server
	server := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Graceful shutdown
	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan

		logger.Info("Shutting down server...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()

		// Shutdown tracer first
		if shutdownTracer != nil {
			if err := shutdownTracer(shutdownCtx); err != nil {
				logger.Error("Failed to shutdown tracer", zap.Error(err))
			}
		}

		if err := server.Shutdown(shutdownCtx); err != nil {
			logger.Error("Server shutdown error", zap.Error(err))
		}
	}()

	// Start server
	logger.Info("Server starting",
		zap.String("port", cfg.Server.Port),
		zap.String("env", cfg.App.Env),
		zap.String("address", fmt.Sprintf("http://localhost:%s", cfg.Server.Port)),
	)

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		logger.Fatal("Server error", zap.Error(err))
	}

	logger.Info("Server stopped gracefully")
}
