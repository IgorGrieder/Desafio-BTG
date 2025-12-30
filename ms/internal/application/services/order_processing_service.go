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
