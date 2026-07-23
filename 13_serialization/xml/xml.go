package main

import (
	"encoding/xml"
	"fmt"
)

type Book struct {
	Title  string
	Author string
}

func main() {
	book := Book{
		"Война и мир",
		"Толстой",
	}

	// Сериаизация
	xmlData2, err := xml.Marshal(book)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Сериализованный XML:")
	fmt.Println(xmlData2)

	// Добавим XML-декларацию для читаемости
	xmlStr := xml.Header + string(xmlData2)
	fmt.Println("Сериализованный XML:")
	fmt.Println(xmlStr)

	// Десериализация
	var book2 Book
	err = xml.Unmarshal(xmlData2, &book2)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Десериализованная книга: %+v\n", book2)

	// Проверка: сериализуем обратно без декларации
	xmlCompact, _ := xml.Marshal(book2)
	fmt.Println("Без декларации:", string(xmlCompact))
}
