// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"
)

// CustomTime сериализуется как строка в формате YYYY-MM-DD
type CustomTime struct {
	time.Time
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return json.Marshal(ct.Format("2006-01-02"))
}

func (ct *CustomTime) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return err
	}
	ct.Time = t
	return nil
}

type Event struct {
	ID   int        `json:"omitempty"`
	Date CustomTime `json:"date"`
}

func main() {
	now := time.Date(2025, 9, 28, 0, 0, 0, 0, time.UTC)
	event := Event{
		ID:   111,
		Date: CustomTime{now},
	}

	// Сериализация
	jsonData, err := json.Marshal(event)
	if err != nil {
		log.Fatal("Ошибка сериализации:", err)
	}
	fmt.Println("Сериализованный JSON:", string(jsonData))

	// Десериализация
	var event2 Event
	err = json.Unmarshal(jsonData, &event2)
	if err != nil {
		log.Fatal("Ошибка десериализации:", err)
	}
	fmt.Printf("Десериализованный объект: ID=%d, Date=%s\n", event2.ID, event2.Date.Format("2006-01-02"))

	// Проверка равенства
	if event.ID == event2.ID && event.Date.Equal(event2.Date.Time) {
		fmt.Println("Кастомная сериализация работает!")
	}
}
