package messaging

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/domain"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/logger"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/ports"
)

type RabbitMQPublisher struct {
	conn     *amqp.Connection
	channel  *amqp.Channel
	exchange string
	queue    string
}

func NewRabbitMQPublisher(url, exchange, queue string) (ports.MessagePublisher, error) {
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// Declare exchange (idempotent operation)
	err = channel.ExchangeDeclare(
		exchange, // name
		"direct", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare queue (idempotent operation)
	_, err = channel.QueueDeclare(
		queue, // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange
	err = channel.QueueBind(
		queue,    // queue name
		queue,    // routing key (same as queue name)
		exchange, // exchange
		false,
		nil,
	)
	if err != nil {
		channel.Close()
		conn.Close()
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	logger.Info("RabbitMQ publisher initialized",
		zap.String("exchange", exchange),
		zap.String("queue", queue),
	)

	return &RabbitMQPublisher{
		conn:     conn,
		channel:  channel,
		exchange: exchange,
		queue:    queue,
	}, nil
}

func (p *RabbitMQPublisher) PublishOrder(ctx context.Context, order *domain.Order) error {
	body, err := json.Marshal(order)
	if err != nil {
		logger.Error("Failed to marshal message",
			zap.Error(err),
			zap.Int64("order_code", order.OrderCode),
		)
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		p.exchange, // exchange
		p.queue,    // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // make message persistent
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		logger.Error("Failed to publish message",
			zap.Error(err),
			zap.Int64("order_code", order.OrderCode),
		)
		return fmt.Errorf("failed to publish message: %v", err)
	}

	logger.Info("Message published successfully",
		zap.Int64("order_code", order.OrderCode),
		zap.Int("customer_code", order.CustomerCode),
		zap.String("exchange", p.exchange),
		zap.String("routing_key", p.queue),
	)

	return nil
}

func (p *RabbitMQPublisher) Close() error {
	logger.Info("Closing RabbitMQ publisher")

	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			logger.Error("Failed to close channel", zap.Error(err))
		}
	}

	if p.conn != nil {
		if err := p.conn.Close(); err != nil {
			logger.Error("Failed to close connection", zap.Error(err))
			return err
		}
	}

	logger.Info("RabbitMQ publisher closed successfully")
	return nil
}
