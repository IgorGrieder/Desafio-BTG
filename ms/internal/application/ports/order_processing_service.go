package ports

import (
	"context"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/domain"
)

// OrderProcessingService defines the interface for order processing business logic
type OrderProcessingService interface {
	// ProcessOrder processes an order message from RabbitMQ and saves to database
	ProcessOrder(ctx context.Context, order *domain.Order) error

	// ValidateOrder validates order data before processing
	ValidateOrder(ctx context.Context, order *domain.Order) error

	// CalculateOrderTotal calculates the total value of an order
	CalculateOrderTotal(order *domain.Order) float64

	// OrderExists checks if an order already exists by code
	OrderExists(ctx context.Context, orderCode int32) (bool, error)

	// GetOrderByCode retrieves an order by its code
	GetOrderByCode(ctx context.Context, orderCode int32) (*domain.Order, error)
}
