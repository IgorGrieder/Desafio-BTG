package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/inbound/consumer"
	db "github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/outbound/database"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/outbound/database/sqlc"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/application/services"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/config"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/logger"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Error("Failed to load configuration", zap.Error(err))
		os.Exit(1)
	}

	// Initialize structured JSON logger
	logger.Init(cfg.App.Env)
	defer logger.Sync()

	logger.Info("Starting Consumer Microservice",
		zap.String("environment", cfg.App.Env),
		zap.String("rabbitmq_host", cfg.RabbitMQ.Host),
		zap.String("rabbitmq_queue", cfg.RabbitMQ.Queue),
		zap.String("database", cfg.Database.DBName),
	)

	// Initialize database connection
	dbConn, err := db.NewDB(ctx, cfg.Database.DSN())
	if err != nil {
		logger.Warn("Failed to connect to database",
			zap.Error(err),
			zap.String("action", "consumer_will_start_without_db"),
		)

		os.Exit(1)
	}

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
		logger.Error("Failed to initialize RabbitMQ consumer",
			zap.Error(err),
			zap.String("rabbitmq_url", cfg.RabbitMQ.Host),
			zap.String("queue", cfg.RabbitMQ.Queue),
		)
		os.Exit(1)
	}
	defer rabbitConsumer.Close()
	logger.Info("RabbitMQ consumer initialized")

	// Start consuming messages
	if err := rabbitConsumer.Start(ctx); err != nil {
		logger.Error("Failed to start consumer", zap.Error(err))
		os.Exit(1)
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

	logger.Info("Consumer stopped gracefully")
}
