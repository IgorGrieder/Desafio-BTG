package services

import (
	"context"
	"fmt"

	db "github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/outbound/database"
	database "github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/adapters/outbound/database/sqlc"
	"github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

// OrderProcessingService handles business logic for processing orders from RabbitMQ
type OrderProcessingService struct {
	queries *db.Store
}

// NewOrderProcessingService creates a new OrderProcessingService with dependency injection
func NewOrderProcessingService(queries *db.Store) *OrderProcessingService {
	return &OrderProcessingService{
		queries: queries,
	}
}

// ProcessOrder processes an order message from RabbitMQ and saves to database
func (s *OrderProcessingService) ProcessOrder(ctx context.Context, order *domain.Order) error {
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

// ValidateOrder validates order data before processing
func (s *OrderProcessingService) ValidateOrder(ctx context.Context, order *domain.Order) error {
	// TODO: Implement validation logic
	// Example:
	// - Check if order already exists
	// - Validate customer code
	// - Validate items
	// - Calculate and verify total
	return nil
}

// CalculateOrderTotal calculates the total value of an order
func (s *OrderProcessingService) CalculateOrderTotal(order *domain.Order) float64 {
	// TODO: Implement calculation logic
	// Example:
	// total := 0.0
	// for _, item := range order.Items {
	//     total += item.Price * float64(item.Quantity)
	// }
	// return total
	return 0.0
}

// OrderExists checks if an order already exists by code
func (s *OrderProcessingService) OrderExists(ctx context.Context, orderCode int32) (bool, error) {
	// TODO: Implement check logic
	// Example:
	// _, err := s.queries.GetOrderByCode(ctx, orderCode)
	// if err != nil {
	//     if err == sql.ErrNoRows {
	//         return false, nil
	//     }
	//     return false, err
	// }
	// return true, nil
	return false, nil
}

// GetOrderByCode retrieves an order by its code
func (s *OrderProcessingService) GetOrderByCode(ctx context.Context, orderCode int32) (*domain.Order, error) {
	// TODO: Implement retrieval logic
	// Example:
	// dbOrder, err := s.queries.GetOrderByCode(ctx, orderCode)
	// if err != nil {
	//     return nil, err
	// }
	// dbItems, err := s.queries.GetOrderItems(ctx, dbOrder.ID)
	// if err != nil {
	//     return nil, err
	// }
	// return convertToOrderDomain(dbOrder, dbItems), nil
	return nil, nil
}
