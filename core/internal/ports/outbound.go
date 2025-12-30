package ports

import (
	"context"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/domain"
)

// MessagePublisher defines the outbound port (driven side) for message broker
// This is what the application needs to publish messages
type MessagePublisher interface {
	// PublishOrder pushes the order to the message broker
	PublishOrder(ctx context.Context, order *domain.Order) error

	// CLoses the pub/sub connection
	Close() error
}
