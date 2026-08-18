package mapper_example

import "time"

//go:generate go run ../cmd/mapper/main.go -type=Order -dto=OrderDTO -output=order_mapper_gen.go order.go

// Order представляет модель заказа
type Order struct {
	ID        int64       `mapper:"id"`
	UserID    int64       `mapper:"user_id"`
	Total     float64     `mapper:"total"`
	Status    string      `mapper:"status"`
	CreatedAt time.Time   `mapper:"created_at"`
	Items     []OrderItem `mapper:"items"`
}

// OrderItem представляет элемент заказа
type OrderItem struct {
	ProductID int64   `mapper:"product_id"`
	Quantity  int     `mapper:"quantity"`
	Price     float64 `mapper:"price"`
}

// OrderDTO представляет DTO для заказа
type OrderDTO struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	Total     float64        `json:"total"`
	Status    string         `json:"status"`
	CreatedAt string         `json:"created_at"`
	Items     []OrderItemDTO `json:"items"`
}

// OrderItemDTO представляет DTO элемента заказа
type OrderItemDTO struct {
	ProductID int64   `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
