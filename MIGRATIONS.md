# Database Migrations with Goose

This project uses [Goose](https://github.com/pressly/goose) for database migrations.

## Installation

Install goose CLI:

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Migration Files Location

- **Core API**: `core/internal/adapters/outbound/database/migrations/`
- **Consumer MS**: `ms/internal/adapters/outbound/database/migrations/`

## Usage

### Run Migrations (Up)

**Core API:**

```bash
cd core
goose -dir internal/adapters/outbound/database/migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=btg_orders sslmode=disable" up
```

**Consumer MS:**

```bash
cd ms
goose -dir internal/adapters/outbound/database/migrations postgres "host=localhost port=5432 user=postgres password=postgres dbname=btg_orders sslmode=disable" up
```

### Rollback Migrations (Down)

```bash
# Rollback one migration
goose -dir internal/adapters/outbound/database/migrations postgres "..." down

# Rollback all migrations
goose -dir internal/adapters/outbound/database/migrations postgres "..." reset
```

### Check Migration Status

```bash
goose -dir internal/adapters/outbound/database/migrations postgres "..." status
```

### Create New Migration

```bash
# In core/
cd core
goose -dir internal/adapters/outbound/database/migrations create add_user_table sql

# In ms/
cd ms
goose -dir internal/adapters/outbound/database/migrations create add_user_table sql
```
