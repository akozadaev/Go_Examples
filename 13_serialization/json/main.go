package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type Person struct {
	Name string `json:"name"`
	Age  string `json:"age"`
}

func main() {
	// Сериализация (marshal)
	p := Person{"Alexey", "45"}
	jsonData, err := json.Marshal(p)
	if err != nil {
		log.Fatal("Ошибка сериализации:", err)
	}
	fmt.Println("Сериализованный JSON:", string(jsonData))

	// Десериализация (unmarshal)
	var p2 Person
	err = json.Unmarshal(jsonData, &p2)
	if err != nil {
		log.Fatal("Ошибка десериализации:", err)
	}
	fmt.Printf("Десериализованный объект: %+v\n", p2)
}
