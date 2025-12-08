package main

import (
	"context"
	"fmt"
	"log"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/outbound/database"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/config"
)

func main() {
	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	fmt.Printf("Starting Order Consumer Microservice\n")
	fmt.Printf("Environment: %s\n", cfg.App.Env)
	fmt.Printf("RabbitMQ Queue: %s\n", cfg.RabbitMQ.Queue)
	fmt.Printf("Database: %s\n", cfg.Database.DBName)

	// Initialize database connection with pgxpool
	db, err := database.NewDB(ctx, cfg.Database.DSN())
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	fmt.Println("Database connection established successfully")

	// TODO: Initialize SQLC queries
	// TODO: Initialize RabbitMQ connection
	// TODO: Start consuming messages

	fmt.Println("Consumer setup complete. Ready to start...")
}
