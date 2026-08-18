package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name  string
	Age   int
	Email string
	City  string
}

type Product struct {
	ID    int
	Name  string
	Price float64
	Stock int
}

// serializeStruct сериализует структуру в map
func serializeStruct(s interface{}) map[string]interface{} {
	v := reflect.ValueOf(s)
	t := v.Type()

	// Проверяем, что передана структура
	if v.Kind() != reflect.Struct {
		return nil
	}

	data := make(map[string]interface{})

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		value := v.Field(i).Interface()
		data[field.Name] = value
	}

	return data
}

// printSerialized выводит сериализованные данные
func printSerialized(s interface{}) {
	data := serializeStruct(s)
	fmt.Printf("Сериализованные данные:\n")
	for key, value := range data {
		fmt.Printf("  %s: %v (%T)\n", key, value, value)
	}
	fmt.Println()
}

// main_serialize демонстрирует сериализацию с использованием рефлексии
func main() {
	fmt.Println("=== Примеры сериализации с рефлексией ===\n")
	/*
		// 1. Сериализация структуры Person
		p := Person{
			Name:  "Alice",
			Age:   30,
			Email: "alice@example.com",
			City:  "New York",
		}

		fmt.Println("--- Структура Person ---")
		fmt.Printf("Исходная структура: %+v\n", p)
		printSerialized(p)*/
	/*
		// 2. Сериализация структуры Product
		prod := Product{
			ID:    1,
			Name:  "Laptop",
			Price: 1299.99,
			Stock: 10,
		}

		fmt.Println("--- Структура Product ---")
		fmt.Printf("Исходная структура: %+v\n", prod)
		printSerialized(prod)*/

	// 3. Демонстрация обработки разных типов
	demonstrateTypeHandling()
}

func demonstrateTypeHandling() {
	fmt.Println("--- Обработка различных типов полей ---")

	type MixedStruct struct {
		StringField string
		IntField    int
		FloatField  float64
		BoolField   bool
		StringSlice []string
		IntMap      map[string]int
	}

	mixed := MixedStruct{
		StringField: "test",
		IntField:    42,
		FloatField:  3.14,
		BoolField:   true,
		StringSlice: []string{"a", "b", "c"},
		IntMap: map[string]int{
			"first":  1,
			"second": 2,
		},
	}

	fmt.Printf("Исходная структура: %+v\n", mixed)
	printSerialized(mixed)

	// Обработка сложных типов
	v := reflect.ValueOf(mixed)
	t := v.Type()

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		fmt.Printf("Поле %s:\n", field.Name)
		fmt.Printf("  Тип: %s\n", field.Type)
		fmt.Printf("  Kind: %s\n", fieldValue.Kind())

		if fieldValue.Kind() == reflect.Slice {
			fmt.Printf("  Длина слайса: %d\n", fieldValue.Len())
			fmt.Printf("  Элементы: ")
			for j := 0; j < fieldValue.Len(); j++ {
				fmt.Printf("%v ", fieldValue.Index(j).Interface())
			}
			fmt.Println()
		} else if fieldValue.Kind() == reflect.Map {
			fmt.Printf("  Элементы map:\n")
			for _, key := range fieldValue.MapKeys() {
				value := fieldValue.MapIndex(key)
				fmt.Printf("    %v: %v\n", key.Interface(), value.Interface())
			}
		}
		fmt.Println()
	}
}
