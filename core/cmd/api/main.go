package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	db "github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/database"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/database/sqlc"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/config"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/server"

	_ "github.com/IgorGrieder/Desafio-BTG/tree/main/core/docs"
)

// @title BTG Pactual Order Processing API
// @version 1.0
// @description API for managing and querying orders
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@btg.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http

func main() {
	ctx := context.Background()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Starting Core API Server on %s:%s\n", cfg.Server.Host, cfg.Server.Port)
	fmt.Printf("Environment: %s\n", cfg.App.Env)
	fmt.Printf("Database: %s\n", cfg.Database.DBName)

	// Initialize database connection
	dbConn, err := db.NewDB(ctx, cfg.Database.DSN())
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Server will start without database connection")
		log.Println("⚠️  Database features will not work. Start PostgreSQL or check .env configuration.")
		dbConn = nil
	} else {
		defer dbConn.Close()
		fmt.Println("Database connection established successfully")
	}

	// Initialize SQLC queries (only if database is connected)
	var queries database.Querier
	if dbConn != nil {
		queries = database.New(dbConn.Pool)
		fmt.Println("✓ SQLC queries initialized")
	} else {
		queries = nil
		fmt.Println("⚠️  SQLC queries not initialized (no database)")
	}

	// Initialize and start HTTP server with dependency injection
	srv := server.NewServer(cfg.Server.Host, cfg.Server.Port, queries)

	// Start server in goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")

	if err := srv.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server stopped gracefully")
}
