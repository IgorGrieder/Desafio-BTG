package ports

import "context"

// OrderService defines the inbound port (driving side)
// This is what HTTP handlers will call
type OrderService interface {
	CreateOrder(ctx context.Context, userID string, amount float64) error
	GetOrder(ctx context.Context, orderID string) (any, error)
	ListOrders(ctx context.Context, userID string) ([]any, error)
}
