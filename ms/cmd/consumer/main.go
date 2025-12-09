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
"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/application/services"
"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/config"
)

func main() {
ctx := context.Background()

// Load configuration
cfg, err := config.Load()
if err != nil {
log.Fatalf("Failed to load configuration: %v", err)
}

fmt.Printf("Starting Consumer Microservice\n")
fmt.Printf("Environment: %s\n", cfg.App.Env)
fmt.Printf("RabbitMQ Host: %s\n", cfg.RabbitMQ.Host)
fmt.Printf("RabbitMQ Queue: %s\n", cfg.RabbitMQ.Queue)
fmt.Printf("Database: %s\n", cfg.Database.DBName)

// Initialize database connection
dbConn, err := db.NewDB(ctx, cfg.Database.DSN())
if err != nil {
log.Printf("Warning: Failed to connect to database: %v", err)
log.Println("Consumer will start without database connection")
log.Println("⚠️  Database features will not work. Start PostgreSQL or check .env configuration.")
dbConn = nil
} else {
defer dbConn.Close()
fmt.Println("✓ Database connection established successfully")
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

// Initialize service with dependency injection
orderProcessingService := services.NewOrderProcessingService(queries)
fmt.Println("✓ OrderProcessingService initialized")

// TODO: Initialize RabbitMQ connection
fmt.Println("⚠️  RabbitMQ connection not implemented yet")

// Start consumer
fmt.Println("\n🚀 Consumer is running...")
fmt.Println("Waiting for messages from RabbitMQ queue:", cfg.RabbitMQ.Queue)
fmt.Println("Press CTRL+C to stop")
fmt.Println("")

// TODO: Start consuming messages
// When message is received, call:
// orderProcessingService.ProcessOrder(ctx, order)

// Keep service running - wait for interrupt signal
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit

fmt.Println("\n⏹️  Shutting down consumer...")

// TODO: Close RabbitMQ connection gracefully

fmt.Println("✓ Consumer stopped gracefully")
}
