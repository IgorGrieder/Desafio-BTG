package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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
	tracer  trace.Tracer
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
		tracer:  otel.Tracer("rabbitmq-consumer"),
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

	// Create a span for the message processing
	ctx, span := c.tracer.Start(ctx, "process_order_message",
		trace.WithSpanKind(trace.SpanKindConsumer),
		trace.WithAttributes(
			attribute.Int64("messaging.delivery_tag", int64(msg.DeliveryTag)),
			attribute.String("messaging.queue", c.queue),
			attribute.Int("messaging.body_size", len(msg.Body)),
		),
	)
	defer span.End()

	// Get trace and span IDs for structured logging
	spanCtx := trace.SpanContextFromContext(ctx)
	traceID := spanCtx.TraceID().String()
	spanID := spanCtx.SpanID().String()

	logger.Info("Message received",
		zap.Uint64("delivery_tag", msg.DeliveryTag),
		zap.Int("size_bytes", len(msg.Body)),
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	)

	var orderMsg *domain.Order
	if err := json.Unmarshal(msg.Body, &orderMsg); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to unmarshal message")
		span.SetAttributes(attribute.String("error.type", "unmarshal_error"))

		logger.Error("Failed to unmarshal message",
			zap.Error(err),
			zap.String("body", string(msg.Body)),
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
		)

		// Reject message without requeue
		if nackErr := msg.Nack(false, false); nackErr != nil {
			logger.Error("Failed to nack message", zap.Error(nackErr))
		}
		return
	}

	// Add order attributes to span
	span.SetAttributes(
		attribute.Int64("order.code", orderMsg.OrderCode),
		attribute.Int("order.customer_code", orderMsg.CustomerCode),
		attribute.Int("order.items_count", len(orderMsg.Items)),
	)

	logger.Info("Processing order",
		zap.Int64("order_code", orderMsg.OrderCode),
		zap.Int("customer_code", orderMsg.CustomerCode),
		zap.String("items: ", fmt.Sprintf("%v", orderMsg.Items)),
		zap.String("created_at", orderMsg.CreatedAt.String()),
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
	)

	// Process the order using the service
	if err := c.service.ProcessOrder(ctx, orderMsg); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to process order")
		span.SetAttributes(attribute.String("error.type", "processing_error"))

		logger.Error("Failed to process order",
			zap.Error(err),
			zap.Int64("order_code", orderMsg.OrderCode),
			zap.Int64("duration_ms", time.Since(startTime).Milliseconds()),
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
		)

		// Reject and requeue the message for retry
		if nackErr := msg.Nack(false, true); nackErr != nil {
			logger.Error("Failed to nack message for requeue", zap.Error(nackErr))
		}
		return
	}

	// Acknowledge the message
	if err := msg.Ack(false); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "failed to ack message")

		logger.Error("Failed to ack message",
			zap.Error(err),
			zap.Int64("order_code", orderMsg.OrderCode),
			zap.String("trace_id", traceID),
			zap.String("span_id", spanID),
		)
		return
	}

	span.SetStatus(codes.Ok, "order processed successfully")

	logger.Info("Order processed successfully",
		zap.Int64("order_code", orderMsg.OrderCode),
		zap.Int64("duration_ms", time.Since(startTime).Milliseconds()),
		zap.String("status", "acknowledged"),
		zap.String("trace_id", traceID),
		zap.String("span_id", spanID),
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
