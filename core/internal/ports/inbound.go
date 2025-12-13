package ports

import (
	"context"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/domain"
)

// OrderService defines the interface for order business logic
type OrderService interface {
	// GetOrderTotal retrieves the total value of an order by code
	GetOrderTotal(ctx context.Context, orderCode int32) (string, error)

	// GetOrderByCode retrieves an order by its code
	GetOrderByCode(ctx context.Context, orderCode int32) (*domain.Order, error)

	// GetOrdersByCustomer retrieves all orders for a customer
	GetOrdersByCustomer(ctx context.Context, customerCode int32) ([]*domain.Order, error)

	// CountOrdersByCustomer counts the number of orders for a customer
	CountOrdersByCustomer(ctx context.Context, customerCode int32) (int64, error)

	// CreateOrder creates a new order with items
	CreateOrder(ctx context.Context, order *domain.Order) error

	// GetOrderItems retrieves all items for an order
	GetOrderItems(ctx context.Context, orderID int64) ([]*domain.OrderItem, error)
}
