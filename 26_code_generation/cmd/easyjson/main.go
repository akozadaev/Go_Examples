package main

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/akozadaev/go-code-generation/easyjson_example"
)

func main() {
	// Демонстрация использования easyjson
	fmt.Println("=== Демонстрация easyjson ===")

	// Создаём тестового пользователя
	user := &easyjson_example.User{
		ID:        1,
		Name:      "Иван Иванов",
		Email:     "ivan@example.com",
		CreatedAt: time.Now(),
		Tags:      []string{"developer", "golang"},
		Profile: &easyjson_example.Profile{
			Bio:      "Go разработчик",
			Settings: map[string]string{"theme": "dark", "lang": "ru"},
			IsActive: true,
			Score:    95.5,
		},
	}

	// Сериализация через easyjson
	jsonEasy, err := user.MarshalJSON()
	if err != nil {
		fmt.Printf("Ошибка easyjson.MarshalJSON: %v\n", err)
	} else {
		fmt.Printf("EasyJSON: %s\n", jsonEasy)
	}

	// Сериализация через стандартный encoding/json для сравнения
	jsonStd, err := json.Marshal(user)
	if err != nil {
		fmt.Printf("Ошибка json.Marshal: %v\n", err)
	} else {
		fmt.Printf("Standard JSON: %s\n", jsonStd)
	}

	// Десериализация через easyjson
	var userUnmarshaled easyjson_example.User
	err = userUnmarshaled.UnmarshalJSON(jsonEasy)
	if err != nil {
		fmt.Printf("Ошибка easyjson.UnmarshalJSON: %v\n", err)
	} else {
		fmt.Printf("Десериализованный пользователь: ID=%d, Name=%s\n", userUnmarshaled.ID, userUnmarshaled.Name)
	}
}

