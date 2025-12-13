package ports

import "context"

// MessagePublisher defines the outbound port (driven side) for message broker
// This is what the application needs to publish messages
type MessagePublisher interface {
	PublishOrder(ctx context.Context, orderID, customerID string, amount float64) error
}
