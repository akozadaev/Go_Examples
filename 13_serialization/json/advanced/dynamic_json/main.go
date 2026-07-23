// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
)

func main() {
	// Динамический JSON (неизвестная структура)
	dynamicJSON := `{
		"type": "user",
		"data": {
			"name": "Alexey",
			"age": 45
		}
	}`

	// Вариант 1: map[string]interface{}
	var raw map[string]interface{}
	if err := json.Unmarshal([]byte(dynamicJSON), &raw); err != nil {
		log.Fatal("Ошибка парсинга в map:", err)
	}
	fmt.Println("Тип события:", raw["type"])
	data := raw["data"].(map[string]interface{})
	fmt.Println("Имя:", data["name"], "Возраст:", data["age"])

	// Вариант 2: Использование json.RawMessage для отложенной обработки
	type Event struct {
		Type string          `json:"type"`
		Data json.RawMessage `json:"data"`
	}

	var event Event
	if err := json.Unmarshal([]byte(dynamicJSON), &event); err != nil {
		log.Fatal("Ошибка парсинга в Event:", err)
	}

	// Теперь обработаем Data отдельно
	var user struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	if err := json.Unmarshal(event.Data, &user); err != nil {
		log.Fatal("Ошибка парсинга Data:", err)
	}

	fmt.Printf("Event.Type = %s, User = %+v\n", event.Type, user)
}
