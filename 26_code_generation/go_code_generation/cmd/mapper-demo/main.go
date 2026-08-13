package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/akozadaev/go-code-generation/mapper_example"
)

func main() {
	fmt.Println("=== Демонстрация генератора мапперов ===")

	// Пример 1: Преобразование User в UserDTO
	fmt.Println("\n1. Преобразование User -> UserDTO:")
	user := &mapper_example.User{
		ID:        1,
		Name:      "Иван Иванов",
		Email:     "ivan@example.com",
		CreatedAt: time.Now(),
		Age:       30,
		IsActive:  true,
		Tags:      []string{"developer", "golang"},
	}

	userDTO := user.ToUserDTO()
	fmt.Printf("Модель: ID=%d, Name=%s, CreatedAt=%v\n", user.ID, user.Name, user.CreatedAt)
	fmt.Printf("DTO: ID=%d, Name=%s, CreatedAt=%s\n", userDTO.ID, userDTO.Name, userDTO.CreatedAt)

	// JSON сериализация DTO
	userJSON, _ := json.Marshal(userDTO)
	fmt.Printf("JSON: %s\n", userJSON)

	// Обратное преобразование DTO -> User
	userFromDTO := mapper_example.UserDTOToUser(userDTO)
	fmt.Printf("Обратно в модель: ID=%d, Name=%s\n", userFromDTO.ID, userFromDTO.Name)

	// Пример 2: Преобразование Product в ProductDTO
	fmt.Println("\n2. Преобразование Product -> ProductDTO:")
	product := &mapper_example.Product{
		ID:          100,
		Name:        "Ноутбук",
		Description: "Мощный ноутбук для разработки",
		Price:       99999.99,
		Stock:       5,
		CreatedAt:   time.Now(),
		CategoryID:  1,
	}

	productDTO := product.ToProductDTO()
	fmt.Printf("Модель: Name=%s, Price=%.2f\n", product.Name, product.Price)
	fmt.Printf("DTO: Name=%s, Price=%.2f, CreatedAt=%s\n", productDTO.Name, productDTO.Price, productDTO.CreatedAt)

	// Пример 3: Преобразование Order в OrderDTO
	fmt.Println("\n3. Преобразование Order -> OrderDTO:")
	order := &mapper_example.Order{
		ID:        1000,
		UserID:    1,
		Total:     199999.98,
		Status:    "pending",
		CreatedAt: time.Now(),
		Items: []mapper_example.OrderItem{
			{ProductID: 100, Quantity: 2, Price: 99999.99},
		},
	}

	orderDTO := order.ToOrderDTO()
	fmt.Printf("Модель: ID=%d, Total=%.2f, Items=%d\n", order.ID, order.Total, len(order.Items))
	fmt.Printf("DTO: ID=%d, Total=%.2f, Items=%d, CreatedAt=%s\n", orderDTO.ID, orderDTO.Total, len(orderDTO.Items), orderDTO.CreatedAt)

	// Пример 4: Преобразование слайсов
	fmt.Println("\n4. Преобразование слайсов:")
	users := []mapper_example.User{
		{ID: 1, Name: "Пользователь 1", Email: "user1@example.com", CreatedAt: time.Now()},
		{ID: 2, Name: "Пользователь 2", Email: "user2@example.com", CreatedAt: time.Now()},
	}

	userDTOs := mapper_example.SliceToUserDTO(users)
	fmt.Printf("Слайс моделей: %d элементов\n", len(users))
	fmt.Printf("Слайс DTO: %d элементов\n", len(userDTOs))
	for i, dto := range userDTOs {
		fmt.Printf("  DTO[%d]: ID=%d, Name=%s\n", i, dto.ID, dto.Name)
	}

	// Обратное преобразование слайса
	usersFromDTO := mapper_example.UserDTOSliceToUser(userDTOs)
	fmt.Printf("Обратно в модели: %d элементов\n", len(usersFromDTO))

	fmt.Println("\n=== Демонстрация завершена ===")
}
