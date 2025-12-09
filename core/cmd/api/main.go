package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	db "github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/database"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/database/sqlc"
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
		logger.Error("Failed to load configuration", slog.String("error", err.Error()))
		os.Exit(1)
	}

	// Initialize structured JSON logger
	logger.Init(cfg.App.Env)

	logger.Info("Starting Core API Server",
		slog.String("host", cfg.Server.Host),
		slog.String("port", cfg.Server.Port),
		slog.String("environment", cfg.App.Env),
		slog.String("database", cfg.Database.DBName),
	)

	// Initialize database connection
	dbConn, err := db.NewDB(ctx, cfg.Database.DSN())
	if err != nil {
		logger.Fatal("Failed to connect to database, shuting down application",
			slog.String("error", err.Error()),
		)
	}

	logger.Info("Database connection established")

	// Initialize SQLC queries
	queries := database.New(dbConn.Pool)

	logger.Info("Repository established")

	// Initialize and start HTTP server with dependency injection
	srv := server.NewServer(cfg.Server.Host, cfg.Server.Port, queries)

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal("Server failed to start",
				slog.String("error", err.Error()))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	if err := srv.Shutdown(); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server stopped gracefully")
}
