package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/inbound/consumer"
	db "github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/outbound/database"
	database "github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/outbound/database/sqlc"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/application/services"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/config"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/infrastructure/telemetry"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// Initialize service with dependency injection
	orderService := services.NewOrderProcessingService(dbStore)

	logger.Info("OrderProcessingService initialized")

	// Initialize RabbitMQ consumer
	rabbitConsumer, err := consumer.NewRabbitMQConsumer(
		cfg.RabbitMQ.URL(),
		cfg.RabbitMQ.Queue,
		orderService,
	)
	if err != nil {
		logger.Fatal("Failed to initialize RabbitMQ consumer",
			zap.Error(err),
			zap.String("rabbitmq_url", cfg.RabbitMQ.Host),
			zap.String("queue", cfg.RabbitMQ.Queue),
		)
	}
	defer rabbitConsumer.Close()

	logger.Info("RabbitMQ consumer initialized")

	// Start consuming messages
	if err := rabbitConsumer.Start(ctx); err != nil {
		logger.Fatal("Failed to start consumer", zap.Error(err))
	}

	logger.Info("Consumer is running and processing messages",
		zap.String("queue", cfg.RabbitMQ.Queue),
		zap.String("status", "active"),
	)

	// Keep service running - wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down consumer gracefully...")

	// Shutdown tracer
	if shutdownTracer != nil {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdownTracer(shutdownCtx); err != nil {
			logger.Error("Failed to shutdown tracer", zap.Error(err))
		}
	}

	logger.Info("Consumer stopped gracefully")
}
