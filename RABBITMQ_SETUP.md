# RabbitMQ Integration Setup

## Overview
Both the Core API and MS Consumer are now integrated with RabbitMQ for asynchronous message processing.

## Architecture

### Core API (Producer)
- **Exchange**: `orders_exchange` (direct)
- **Queue**: `orders`
- **Routing Key**: `orders`
- Publishes order messages when orders are created

### MS Consumer
- **Queue**: `orders`
- Consumes messages from the same queue
- Processes orders asynchronously
- Acknowledges messages after successful processing

## Configuration

### Core API (.env)
```bash
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_EXCHANGE=orders_exchange
RABBITMQ_QUEUE=orders
```

### MS Consumer (.env)
```bash
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
RABBITMQ_QUEUE=orders
```

## Message Format
```json
{
  "order_id": "string",
  "customer_id": "string",
  "amount": 123.45,
  "status": "pending",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

## Components Created

### Core API
- `/internal/adapters/outbound/messaging/message-publisher.go` - RabbitMQ publisher adapter
- `/internal/ports/outbound.go` - MessagePublisher interface (outbound port)
- Updated service to inject MessagePublisher dependency

### MS Consumer
- `/internal/adapters/inbound/consumer/rabbitmq.go` - RabbitMQ consumer adapter
- Implements message acknowledgment and retry logic
- QoS set to 1 (process one message at a time)

## How It Works

1. **Core API** receives HTTP request to create order
2. **Core Service** processes the order
3. **MessagePublisher** publishes message to RabbitMQ exchange
4. **RabbitMQ** routes message to `orders` queue
5. **MS Consumer** receives and processes the message
6. **MS Service** handles business logic
7. **MS Consumer** acknowledges the message

## Running

### Start RabbitMQ
```bash
docker-compose up -d rabbitmq
```

### Start Core API
```bash
cd core
go run cmd/api/main.go
```

### Start MS Consumer
```bash
cd ms
go run cmd/consumer/main.go
```

## Hexagonal Architecture Ports

### Outbound Port (Driven Side)
- **Interface**: `ports.MessagePublisher`
- **Implementation**: `messaging.RabbitMQPublisher`
- **Purpose**: What the application needs to publish messages

### Inbound Port (Driving Side) - MS
- **Interface**: `ports.MessageProcessor`
- **Implementation**: `consumer.RabbitMQConsumer`
- **Purpose**: How external systems trigger the application

## Error Handling

### Publisher
- Connection failures logged and returned
- Marshal errors logged with order details
- Publish failures logged with order context

### Consumer
- Invalid JSON messages: Nack without requeue
- Processing failures: Nack with requeue for retry
- Successful processing: Ack to remove from queue

## Logging
Both services use structured JSON logging with Zap:
- Connection events
- Message publishing/consuming
- Processing duration
- Error details with context
