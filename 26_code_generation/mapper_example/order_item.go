package mapper_example

//go:generate go run ../cmd/mapper/main.go -type=OrderItem -dto=OrderItemDTO -output=order_item_mapper_gen.go order.go

// Этот файл содержит только директиву go:generate для генерации маппера OrderItem
// Структуры определены в order.go
