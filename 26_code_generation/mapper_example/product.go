package mapper_example

import "time"

//go:generate go run ../cmd/mapper/main.go -type=Product -dto=ProductDTO -output=product_mapper_gen.go product.go

// Product представляет модель товара
type Product struct {
	ID          int64     `mapper:"id"`
	Name        string    `mapper:"name"`
	Description string    `mapper:"description"`
	Price       float64   `mapper:"price"`
	Stock       int       `mapper:"stock"`
	CreatedAt   time.Time `mapper:"created_at"`
	CategoryID  int64     `mapper:"category_id"`
}

// ProductDTO представляет DTO для товара
type ProductDTO struct {
	ID          int64   `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
	CreatedAt   string  `json:"created_at"`
	CategoryID  int64   `json:"category_id"`
}
