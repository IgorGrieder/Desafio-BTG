package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	db "github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/outbound/database"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/outbound/database/sqlc"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/application/services"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/config"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/logger"
)

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize structured JSON logger
	logger.Init(cfg.App.Env)

	logger.Info("Starting Consumer Microservice",
		"environment", cfg.App.Env,
		"rabbitmq_host", cfg.RabbitMQ.Host,
		"rabbitmq_queue", cfg.RabbitMQ.Queue,
		"database", cfg.Database.DBName,
	)

	// Initialize database connection
	dbConn, err := db.NewDB(ctx, cfg.Database.DSN())
	if err != nil {
		logger.Warn("Failed to connect to database",
			"error", err,
			"action", "consumer_will_start_without_db",
		)
		logger.Warn("Database features will not work",
			"suggestion", "start_postgresql_or_check_env_config",
		)
		dbConn = nil
	} else {
		defer dbConn.Close()
		logger.Info("Database connection established")
	}

	// Initialize SQLC queries (only if database is connected)
	var queries database.Querier
	if dbConn != nil {
		queries = database.New(dbConn.Pool)
		logger.Info("SQLC queries initialized")
	} else {
		queries = nil
		logger.Warn("SQLC queries not initialized", "reason", "no_database_connection")
	}

	// Initialize service with dependency injection
	_ = services.NewOrderProcessingService(queries)
	logger.Info("OrderProcessingService initialized")

	// TODO: Initialize RabbitMQ connection
	logger.Warn("RabbitMQ connection not implemented yet", "status", "pending")

	// Start consumer
	logger.Info("Consumer is running",
		"queue", cfg.RabbitMQ.Queue,
		"status", "waiting_for_messages",
	)

	// TODO: Start consuming messages
	// When message is received, call:
	// logger.Info("Message received", "order_id", order.ID)
	// orderProcessingService.ProcessOrder(ctx, order)

	// Keep service running - wait for interrupt signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down consumer gracefully...")

	// TODO: Close RabbitMQ connection gracefully

	logger.Info("Consumer stopped gracefully")
}
