package main

import (
"context"
"log"

"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/outbound/database/sqlc"
"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/application/services"
)

type Consumer struct {
ctx                     context.Context
orderProcessingService  *services.OrderProcessingService
}

func NewConsumer(ctx context.Context, queries database.Querier) *Consumer {
// Initialize service with dependency injection
orderProcessingService := services.NewOrderProcessingService(queries)

return &Consumer{
ctx:                    ctx,
orderProcessingService: orderProcessingService,
}
}

func (c *Consumer) Start() error {
log.Println("Consumer started")
log.Println("✓ OrderProcessingService initialized with SQLC queries")
log.Println("Listening to RabbitMQ queue...")

// TODO: Start consuming from RabbitMQ
// TODO: On message received, call c.orderProcessingService.ProcessOrder(ctx, order)

// Block forever (until context is cancelled)
<-c.ctx.Done()
return nil
}

func (c *Consumer) Shutdown() error {
log.Println("Consumer shutting down...")
// TODO: Close RabbitMQ connection
// TODO: Close database connection
return nil
}
