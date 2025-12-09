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
		slog.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// Initialize structured JSON logger
	logger.Init(cfg.App.Env)

	logger.Info("Starting Core API Server",
		"host", cfg.Server.Host,
		"port", cfg.Server.Port,
		"environment", cfg.App.Env,
		"database", cfg.Database.DBName,
	)

	// Initialize database connection
	dbConn, err := db.NewDB(ctx, cfg.Database.DSN())
	if err != nil {
		logger.Warn("Failed to connect to database",
			"error", err,
			"action", "server_will_start_without_db",
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

	// Initialize and start HTTP server with dependency injection
	srv := server.NewServer(cfg.Server.Host, cfg.Server.Port, queries)

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			logger.Fatal("Server failed to start", "error", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server gracefully...")

	if err := srv.Shutdown(); err != nil {
		logger.Fatal("Server forced to shutdown", "error", err)
	}

	logger.Info("Server stopped gracefully")
}
