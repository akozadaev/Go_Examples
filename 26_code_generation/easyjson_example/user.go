package easyjson_example

//go:generate easyjson -all user.go

import (
	"time"
)

// User представляет пользователя системы
//
//easyjson:json
type User struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Tags      []string  `json:"tags"`
	Profile   *Profile  `json:"profile,omitempty"`
}

// Profile представляет профиль пользователя
//
//easyjson:json
type Profile struct {
	Bio      string            `json:"bio"`
	Settings map[string]string `json:"settings"`
	IsActive bool              `json:"is_active"`
	Score    float64           `json:"score"`
}
