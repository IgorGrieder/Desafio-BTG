package domain

import (
	"time"
)

type Order struct {
	CustomerCode int
	OrderCode    int64
	Items        []OrderItem
	CreatedAt    time.Time
}

type OrderItem struct {
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
