package ports

import "github.com/IgorGrieder/Desafio-BTG/tree/main/ms/internal/domain"

type OrderRepository interface {
	Save(order *domain.Order) error
	FindByID(id int64) (*domain.Order, error)
	FindByCustomerCode(customerCode int) ([]*domain.Order, error)
	GetTotalByOrderCode(orderCode int) (float64, error)
	CountOrdersByCustomer(customerCode int) (int, error)
}
