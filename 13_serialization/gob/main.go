package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
)

type Book struct {
	Title  string
	Author string
}

func main() {
	book := Book{
		Title:  "Go Programming Language",
		Author: "Gopher",
	}

	// Буфер для хранения данных
	var buffer bytes.Buffer

	enc := gob.NewEncoder(&buffer)
	err := enc.Encode(book)
	if err != nil {
		log.Fatal("Ошибка gob.Encode:", err)
	}

	// Десериализация
	var book2 Book

	dec := gob.NewDecoder(&buffer)

	err = dec.Decode(&book2)
	if err != nil {
		log.Fatal("Ошибка gob.Decode:", err)
	}
	fmt.Printf("Исходный: %+v\n", book)
	fmt.Printf("Восстановленный: %+v\n", book2)

	// Проверим, что данные идентичны
	if book.Title == book2.Title && book.Author == book2.Author {
		fmt.Println("Gob сериализация/десериализация прошла успешно!")
	}
}
