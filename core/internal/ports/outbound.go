package ports

import "context"

// OrderRepository defines the outbound port (driven side) for database
// This is what the application needs from the database adapter
type OrderRepository interface {
	CreateOrder(ctx context.Context, userID string, amount float64) error
	GetOrder(ctx context.Context, orderID string) (any, error)
	ListOrders(ctx context.Context, userID string) ([]any, error)
}

// MessagePublisher defines the outbound port (driven side) for message broker
// This is what the application needs to publish messages
type MessagePublisher interface {
	PublishOrder(ctx context.Context, orderID string, data any) error
}
