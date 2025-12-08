package domain

import "time"

type Order struct {
	ID             int64
	Code           int
	CustomerCode   int
	Items          []OrderItem
	TotalValue     float64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type OrderItem struct {
	ID       int64
	OrderID  int64
	Product  string
	Quantity int
	Price    float64
}

func (o *Order) CalculateTotal() float64 {
	var total float64
	for _, item := range o.Items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}
