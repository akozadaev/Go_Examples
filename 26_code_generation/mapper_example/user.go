package mapper_example

import "time"

//go:generate go run ../cmd/mapper/main.go -type=User -dto=UserDTO -output=user_mapper_gen.go user.go

// User представляет модель пользователя в базе данных
type User struct {
	ID        int64     `mapper:"id"`
	Name      string    `mapper:"name"`
	Email     string    `mapper:"email"`
	CreatedAt time.Time `mapper:"created_at"`
	Age       int       `mapper:"age"`
	IsActive  bool      `mapper:"is_active"`
	Tags      []string  `mapper:"tags"`
}

// UserDTO представляет DTO для передачи данных через API
type UserDTO struct {
	ID        int64    `json:"id"`
	Name      string   `json:"name"`
	Email     string   `json:"email"`
	CreatedAt string   `json:"created_at"` // Формат: RFC3339
	Age       int      `json:"age"`
	IsActive  bool     `json:"is_active"`
	Tags      []string `json:"tags"`
}
