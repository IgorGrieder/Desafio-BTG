package main

import (
	"context"
	"log"
)

type Consumer struct {
	ctx context.Context
}

func NewConsumer(ctx context.Context) *Consumer {
	return &Consumer{
		ctx: ctx,
	}
}

func (c *Consumer) Start() error {
	log.Println("Consumer started")
	log.Println("Listening to RabbitMQ queue...")

	// TODO: Start consuming from RabbitMQ
	// TODO: Process messages and save to database

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
