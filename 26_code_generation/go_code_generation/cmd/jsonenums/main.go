package main

import (
	"encoding/json"
	"fmt"

	"github.com/akozadaev/go-code-generation/jsonenums_example"
)

func main() {
	// Демонстрация использования jsonenums
	priorities := []jsonenums_example.Priority{
		jsonenums_example.PriorityLow,
		jsonenums_example.PriorityMedium,
		jsonenums_example.PriorityHigh,
		jsonenums_example.PriorityCritical,
	}

	fmt.Println("=== Демонстрация jsonenums ===")

	// Тест сериализации (MarshalJSON)
	for _, priority := range priorities {
		jsonData, err := json.Marshal(priority)
		if err != nil {
			fmt.Printf("Ошибка при сериализации %v: %v\n", priority, err)
			continue
		}
		fmt.Printf("Priority: %d, JSON: %s\n", priority, jsonData)
	}

	// Тест десериализации (UnmarshalJSON)
	jsonStrings := []string{`"PriorityLow"`, `"PriorityHigh"`, `"invalid"`}
	for _, jsonStr := range jsonStrings {
		var priority jsonenums_example.Priority
		err := json.Unmarshal([]byte(jsonStr), &priority)
		if err != nil {
			fmt.Printf("Ошибка при десериализации %s: %v\n", jsonStr, err)
		} else {
			fmt.Printf("JSON: %s, Priority: %d\n", jsonStr, priority)
		}
	}
}

