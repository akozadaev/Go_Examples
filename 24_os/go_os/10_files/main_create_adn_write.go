package main

import (
	"fmt"
	"os"
)

func main() {
	// Создание или открытие файла для записи
	file, err := os.Create("example.txt") // os.CreateTemp()
	//file, err := os.Create("example.txt") // os.CreateTemp()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Запись данных в файл
	data := []byte("Hello, World!")
	n, err := file.Write(data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("Записано %d байт\n", n)
}
