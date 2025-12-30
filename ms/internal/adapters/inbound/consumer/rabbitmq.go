package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/application/services"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/domain"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

type RabbitMQConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	queue   string
	service *services.OrderProcessingService
}

type OrderMessage struct {
	OrderID    string  `json:"order_id"`
	CustomerID string  `json:"customer_id"`
	Amount     float64 `json:"amount"`
	Status     string  `json:"status"`
}

func NewRabbitMQConsumer(url, queueName string, service *services.OrderProcessingService) (*RabbitMQConsumer, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare queue (idempotent operation)
	_, err = channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Set QoS - process one message at a time
	err = channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to set QoS: %w", err)
	}

	return &RabbitMQConsumer{
		conn:    conn,
		channel: channel,
		queue:   queueName,
		service: service,
	}, nil
}

func (c *RabbitMQConsumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queue, // queue
		"",      // consumer
		false,   // auto-ack (disabled - we'll ack manually)
		false,   // exclusive
		false,   // no-local
		false,   // no-wait
		nil,     // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	logger.Info("Consumer started successfully",
		zap.String("queue", c.queue),
		zap.String("status", "waiting_for_messages"),
	)

	go func() {
		for {
			select {
			case <-ctx.Done():
				logger.Info("Consumer context cancelled", zap.String("reason", "shutdown"))
				return
			case msg, ok := <-msgs:
				if !ok {
					logger.Warn("Message channel closed")
					return
				}

				c.processMessage(ctx, msg)
			}
		}
	}()

	return nil
}

func (c *RabbitMQConsumer) processMessage(ctx context.Context, msg amqp.Delivery) {
	startTime := time.Now()

	logger.Info("Message received",
		zap.Uint64("delivery_tag", msg.DeliveryTag),
		zap.Int("size_bytes", len(msg.Body)),
	)

	var orderMsg *domain.Order
	if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
		logger.Error("Failed to unmarshal message",
			zap.Error(err),
			zap.String("body", string(msg.Body)),
		)

		// Reject message without requeue
		if nackErr := msg.Nack(false, false); nackErr != nil {
			logger.Error("Failed to nack message", zap.Error(nackErr))
		}
		return
	}

	logger.Info("Processing order",
		zap.Int64("order code", orderMsg.OrderCode),
		zap.Int("customer code", orderMsg.CustomerCode),
		zap.String("items: ", fmt.Sprintf("%v", orderMsg.Items)),
		zap.String("created at", orderMsg.CreatedAt.String()),
	)

	// Process the order using the service
	if err := c.service.ProcessOrder(ctx, orderMsg); err != nil {
		logger.Error("Failed to process order",
			zap.Error(err),
			zap.Int64("order code", orderMsg.OrderCode),
			zap.Int64("duration_ms", time.Since(startTime).Milliseconds()),
		)

		// Reject and requeue the message for retry
		if nackErr := msg.Nack(false, true); nackErr != nil {
			logger.Error("Failed to nack message for requeue", zap.Error(nackErr))
		}
		return
	}

	// Acknowledge the message
	if err := msg.Ack(false); err != nil {
		logger.Error("Failed to ack message",
			zap.Error(err),
			zap.Int64("order code", orderMsg.OrderCode),
		)
		return
	}

	logger.Info("Order processed successfully",
		zap.Int64("order code", orderMsg.OrderCode),
		zap.Int64("duration_ms", time.Since(startTime).Milliseconds()),
		zap.String("status", "acknowledged"),
	)
}

func (c *RabbitMQConsumer) Close() error {
	logger.Info("Closing RabbitMQ consumer")

	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			logger.Error("Failed to close channel", zap.Error(err))
		}
	}

	if c.conn != nil {
		if err := c.conn.Close(); err != nil {
			logger.Error("Failed to close connection", zap.Error(err))
			return err
		}
	}

	logger.Info("RabbitMQ consumer closed successfully")
	return nil
}
