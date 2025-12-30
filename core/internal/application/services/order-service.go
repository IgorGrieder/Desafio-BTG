package services

import (
	"context"
	"fmt"

	db "github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/database"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/adapters/outbound/database/sqlc"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/IgorGrieder/Desafio-BTG/tree/main/core/internal/ports"
)

// Handles DB conns and queries
type OrderServiceQuerier struct {
	*database.Queries
}

// OrderService handles business logic for orders
type OrderService struct {
	queries          *db.Store
	messagePublisher ports.MessagePublisher
}

// NewOrderService creates a new OrderService with dependency injection
func NewOrderService(queries *db.Store, messagePublisher ports.MessagePublisher) ports.OrderService {
	return &OrderService{
		queries:          queries,
		messagePublisher: messagePublisher,
	}
}

// GetOrderTotal retrieves the total value of an order by code
func (s *OrderService) GetOrderTotal(ctx context.Context, orderCode int32) (string, error) {
	// TODO: Implement business logic
	// Example:
	// total, err := s.queries.GetTotalByOrderCode(ctx, orderCode)
	// if err != nil {
	//     return "", err
	// }
	// return total, nil
	tx, err := s.queries.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return "", err
	}
	s.queries.WithTx(tx)
	return "", nil
}

// GetOrderByCode retrieves an order by its code
func (s *OrderService) GetOrderByCode(ctx context.Context, orderCode int32) (*domain.Order, error) {
	// TODO: Implement business logic
	// Example:
	// dbOrder, err := s.queries.GetOrderByCode(ctx, orderCode)
	// if err != nil {
	//     return nil, err
	// }
	// return convertToOrderDomain(dbOrder), nil
	return nil, nil
}

// GetOrdersByCustomer retrieves all orders for a customer
func (s *OrderService) GetOrdersByCustomer(ctx context.Context, customerCode int32) ([]*domain.Order, error) {
	// TODO: Implement business logic
	// Example:
	// dbOrders, err := s.queries.GetOrdersByCustomerCode(ctx, customerCode)
	// if err != nil {
	//     return nil, err
	// }
	// return convertToOrdersDomain(dbOrders), nil
	return nil, nil
}

// CountOrdersByCustomer counts the number of orders for a customer
func (s *OrderService) CountOrdersByCustomer(ctx context.Context, customerCode int32) (int64, error) {
	// TODO: Implement business logic
	// Example:
	// count, err := s.queries.CountOrdersByCustomer(ctx, customerCode)
	// if err != nil {
	//     return 0, err
	// }
	// return count, nil
	return 0, nil
}

// CreateOrder creates a new order with items
func (s *OrderService) CreateOrder(ctx context.Context, order *domain.Order) error {
	tx, err := s.queries.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return fmt.Errorf("error opening transaction %v", err)
	}
	defer tx.Commit(ctx)

	queries := s.queries.WithTx(tx)

	args := database.CreateOrderParams{
		Code:         int32(order.OrderCode),
		CustomerCode: int32(order.CustomerCode),
	}

	orderCreated, err := queries.CreateOrder(ctx, args)
	if err != nil {
		tx.Rollback(ctx)
		return fmt.Errorf("error creating order %v", err)
	}

	orderId := int64(orderCreated.Code)

	for _, item := range order.Items {
		var p pgtype.Numeric

		err := p.Scan(item.Price)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}

		args := database.CreateOrderItemParams{
			OrderID:  orderId,
			Product:  item.Product,
			Price:    p,
			Quantity: int32(item.Quantity),
		}

		_, err = queries.CreateOrderItem(ctx, args)
		if err != nil {
			tx.Rollback(ctx)
			return fmt.Errorf("error creating order item %v : %v", item, err)
		}
	}

	return nil
}

// GetOrderItems retrieves all items for an order
func (s *OrderService) GetOrderItems(ctx context.Context, orderID int64) ([]*domain.OrderItem, error) {
	// TODO: Implement business logic
	// Example:
	// dbItems, err := s.queries.GetOrderItems(ctx, orderID)
	// if err != nil {
	//     return nil, err
	// }
	// return convertToOrderItemsDomain(dbItems), nil
	return nil, nil
}
