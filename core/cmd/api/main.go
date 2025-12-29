package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	db "github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/database"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/database/sqlc"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/messaging"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/config"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/logger"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/server"

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
		logger.Error("Failed to load configuration", zap.Error(err))
		os.Exit(1)
	}

	// Initialize structured JSON logger
	logger.Init(cfg.App.Env)
	defer logger.Sync()

	logger.Info("Starting Core API Server",
		zap.String("host", cfg.Server.Host),
		zap.String("port", cfg.Server.Port),
		zap.String("environment", cfg.App.Env),
		zap.String("database", cfg.Database.DBName),
	)

	// Initialize database connection
	dbConn, err := db.NewDB(ctx, cfg.Database.DSN())
	if err != nil {
		logger.Fatal("Failed to connect to database, shuting down application",
			zap.Error(err),
		)
	}

	logger.Info("Database connection established")

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

	// Initialize and start HTTP server with dependency injection
	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	srv := server.NewServer(cfg.Server.Host, cfg.Server.Port, dbStore, publisher)

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			logger.Error("Server failed to start",
				zap.Error(err))
			quit <- syscall.SIGTERM // Trigger shutdown
		}
	}()

	<-quit

	if err := srv.Shutdown(); err != nil {
		logger.Error("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server stopped gracefully")
}
