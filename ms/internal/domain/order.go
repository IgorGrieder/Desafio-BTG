package domain

import "time"

type Order struct {
	CustomerCode int         `json:"customerCode"`
	OrderCode    int64       `json:"orderCode"`
	Items        []OrderItem `json:"items"`
	CreatedAt    time.Time   `json:"createdAt"`
}

type OrderItem struct {
	Product  string  `json:"product"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func (o *Order) CalculateTotal() float64 {
	var total float64
	for _, item := range o.Items {
		total += float64(item.Quantity) * item.Price
	}
	return total
}
