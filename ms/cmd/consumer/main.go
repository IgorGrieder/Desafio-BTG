package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	db "github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/outbound/database"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/outbound/database/sqlc"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/config"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Starting Order Consumer Microservice\n")
	fmt.Printf("Environment: %s\n", cfg.App.Env)
	fmt.Printf("RabbitMQ Queue: %s\n", cfg.RabbitMQ.Queue)
	fmt.Printf("Database: %s\n", cfg.Database.DBName)

	// Initialize database connection
	dbConn, err := db.NewDB(ctx, cfg.Database.DSN())
	if err != nil {
		log.Printf("Warning: Failed to connect to database: %v", err)
		log.Println("Consumer will start without database connection")
	} else {
		defer dbConn.Close()
		fmt.Println("Database connection established successfully")
	}

	// Initialize SQLC queries
	queries := database.New(dbConn.Pool)
	fmt.Println("✓ SQLC queries initialized")

	// TODO: Initialize RabbitMQ connection

	// Initialize and start consumer with dependency injection
	consumer := NewConsumer(ctx, queries)

	// Start consumer in goroutine
	go func() {
		if err := consumer.Start(); err != nil {
			log.Fatalf("Consumer failed to start: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down consumer...")
	cancel() // Cancel context to stop consumer

	if err := consumer.Shutdown(); err != nil {
		log.Fatalf("Consumer forced to shutdown: %v", err)
	}

	fmt.Println("Consumer stopped gracefully")
}
