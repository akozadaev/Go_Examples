// main.go
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

type Settings struct {
	Debug   bool     `json:"debug"`
	Servers []string `json:"servers"`
}

func main() {
	settings := Settings{
		Debug:   true,
		Servers: []string{"server1", "server2"},
	}

	// Запись в файл
	file, err := os.Create("config.json")
	if err != nil {
		log.Fatal("Не удалось создать файл:", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ") // красивый вывод (не в одну строку)
	if err = encoder.Encode(settings); err != nil {
		log.Fatal("Ошибка записи в файл:", err)
	}

	fmt.Printf("Путь к файлу %s\n", file.Name())
	fmt.Println("Конфигурация сохранена в config.json")

	// Чтение из файла
	file2, err := os.Open("config.json")
	if err != nil {
		log.Fatal("Не удалось открыть файл:", err)
	}
	defer file2.Close()

	var loaded Settings
	if err = json.NewDecoder(file2).Decode(&loaded); err != nil {
		log.Fatal("Ошибка чтения из файла:", err)
	}

	fmt.Printf("Загружено из файла: %+v\n", loaded)
}
