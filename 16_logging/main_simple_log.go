package main

import "log"

func main() {
	log.Println("Это сообщение будет залогировано")

	err := someFunction()
	if err != nil {
		log.Printf("Произошла ошибка: %v", err)
	}
}

func someFunction() error {
	// Имитация ошибки
	return nil
}
