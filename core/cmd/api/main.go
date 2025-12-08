package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/database"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/config"

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
	db, err := database.NewDB(ctx, cfg.Database.DSN())
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Server will start without database connection")
	} else {
		defer db.Close()
		fmt.Println("Database connection established successfully")
	}

	// Initialize and start HTTP server
	server := NewServer(cfg.Server.Host, cfg.Server.Port)

	// Start server in goroutine
	go func() {
		if err := server.Start(); err != nil {
			log.Fatalf("Server failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down server...")

	if err := server.Shutdown(); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	fmt.Println("Server stopped gracefully")
}

