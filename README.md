# BTG Pactual - Order Processing System

This project implements an order processing system using hexagonal architecture in Go.

## Architecture

The project is divided into two main services:

- **core**: REST API service for querying order information
- **ms**: Microservice consumer that processes orders from RabbitMQ

### Hexagonal Architecture Structure

```
├── core/                          # REST API Service
│   ├── cmd/
│   │   └── api/
│   │       └── main.go            # Application entry point
│   ├── internal/
│   │   ├── config/                # Configuration (private)
│   │   ├── domain/                # Business entities
│   │   ├── application/
│   │   │   ├── ports/             # Interfaces (repositories, services)
│   │   │   └── services/          # Business logic
│   │   └── adapters/
│   │       ├── inbound/
│   │       │   └── http/          # HTTP handlers
│   │       └── outbound/
│   │           ├── database/      # Database implementations
│   │           └── messaging/     # Message queue implementations
│
└── ms/                            # Consumer Microservice
    ├── cmd/consumer/              # Application entry point
    ├── internal/
    │   ├── config/                # Configuration (private)
    │   ├── domain/                # Business entities
    │   ├── application/
    │   │   ├── ports/             # Interfaces
    │   │   └── services/          # Business logic
    │   └── adapters/
    │       ├── inbound/
    │       │   └── messaging/     # RabbitMQ consumer
    │       └── outbound/
    │           └── database/      # Database implementations
```

## Prerequisites

- Go 1.25 or higher
- Docker and Docker Compose
- PostgreSQL (via Docker)
- RabbitMQ (via Docker)

## Setup

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd btg
   ```

2. **Start infrastructure services**

   ```bash
   docker-compose up -d
   ```

3. **Configure environment variables**

   For the core service:

   ```bash
   cd core
   cp .env.example .env
   # Edit .env with your configurations
   ```

   For the microservice:

   ```bash
   cd ms
   cp .env.example .env
   # Edit .env with your configurations
   ```

4. **Install dependencies**

   ```bash
   # For core service
   cd core
   go mod download

   # For microservice
   cd ../ms
   go mod download
   ```

## Running the Services

### Core API Service

```bash
cd core
go run ./cmd/api
```

The API will be available at `http://localhost:8080`

### Consumer Microservice

```bash
cd ms
go run main.go
```

## Environment Variables

### Core Service (.env)

- `PORT`: HTTP server port (default: 8080)
- `HOST`: HTTP server host (default: localhost)
- `DB_HOST`: PostgreSQL host
- `DB_PORT`: PostgreSQL port
- `DB_USER`: Database user
- `DB_PASSWORD`: Database password
- `DB_NAME`: Database name
- `DB_SSL_MODE`: SSL mode for database connection
- `APP_ENV`: Application environment (development/production)
- `LOG_LEVEL`: Logging level

### Microservice (.env)

- `RABBITMQ_HOST`: RabbitMQ host
- `RABBITMQ_PORT`: RabbitMQ port
- `RABBITMQ_USER`: RabbitMQ user
- `RABBITMQ_PASSWORD`: RabbitMQ password
- `RABBITMQ_QUEUE`: Queue name to consume
- `DB_*`: Same database configuration as core service

## Docker Services

- **PostgreSQL**: `localhost:5432`
- **RabbitMQ**:
  - AMQP: `localhost:5672`
  - Management UI: `http://localhost:15672` (guest/guest)
- **Jaeger** (Tracing): `http://localhost:16686`
- **Prometheus** (Metrics): `http://localhost:9090`
- **Grafana** (Dashboards): `http://localhost:3000` (admin/admin)

## Observability

The project includes full observability with tracing, metrics, and logging.

### OpenTelemetry Tracing

Distributed tracing is enabled via OpenTelemetry with Jaeger as the backend.

**Configuration:**
```bash
OTEL_ENABLED=true
OTEL_EXPORTER_OTLP_ENDPOINT=http://localhost:4318
```

**Features:**
- Automatic HTTP request tracing with `otelhttp`
- Custom spans for message processing in the consumer
- W3C TraceContext propagation for distributed tracing
- Trace ID and Span ID included in all log messages

**Access Jaeger UI:** `http://localhost:16686`

### Prometheus Metrics

The core API exposes Prometheus metrics for request monitoring.

**Metrics exposed:**
- `http_requests_total` - Counter of total HTTP requests (by method, path, status)
- `http_request_duration_seconds` - Histogram of request durations (by method, path)

**Metrics endpoint:** `GET /metrics`

**Access Prometheus UI:** `http://localhost:9090`

### Grafana Dashboards

Pre-configured Grafana dashboards for API monitoring.

**Access Grafana:** `http://localhost:3000` (admin/admin)

**Pre-configured datasources:**
- Prometheus (metrics)
- Jaeger (traces)

**Included dashboard:**
- API Dashboard with request rate, p95 latency, and error rate panels

### Structured Logging

Both services use structured JSON logging with zap.

**Log fields:**
- `timestamp` - ISO 8601 timestamp
- `level` - Log level (info, warn, error)
- `message` - Log message
- `trace_id` - OpenTelemetry trace ID
- `span_id` - OpenTelemetry span ID
- Additional context fields per log

### Health Endpoints

**Core API:**
- `GET /health` - Health check endpoint
- `GET /metrics` - Prometheus metrics

## API Endpoints

- `GET /orders/:code/total` - Get total value of an order
- `GET /customers/:code/orders/count` - Get number of orders by customer
- `GET /customers/:code/orders` - Get list of orders by customer

## Order Message Format

```json
{
  "codigoPedido": 1001,
  "codigoCliente": 1,
  "itens": [
    {
      "produto": "lápis",
      "quantidade": 100,
      "preco": 1.1
    },
    {
      "produto": "caderno",
      "quantidade": 10,
      "preco": 1.0
    }
  ]
}
```

## Architecture Diagram

📊 **Visual Diagram**: Open `architecture.drawio` with [draw.io](https://app.diagrams.net/) for a complete visual representation of the hexagonal architecture.


## License

MIT
