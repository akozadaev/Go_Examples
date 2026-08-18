package main

import (
	"fmt"
	"reflect"
)

// main_valueof демонстрирует использование reflect.Value и ValueOf
func main() {
	fmt.Println("=== Примеры использования reflect.ValueOf ===")

	// 1. Получение значения переменной
	x := 42
	v := reflect.ValueOf(x)
	fmt.Printf("Значение переменной x через рефлексию: %v\n", v.Interface())
	fmt.Printf("Тип значения: %s\n", v.Type())
	fmt.Printf("Kind: %s\n", v.Kind())

	// 2. Изменение значения переменной
	modifyValue()

	// 3. Получение значения из указателя
	getValueFromPointer()
}

func modifyValue() {
	fmt.Println("\n--- Изменение значения переменной ---")

	x := 42
	v := reflect.ValueOf(&x) // Получаем указатель на переменную

	// Элемент указывает на фактическое значение
	fmt.Printf("Значение до изменения: %v\n", x)

	// Изменяем значение через рефлексию
	v.Elem().SetInt(24)
	fmt.Printf("Значение после изменения: %v\n", x)

	// Попытка изменить не settable значение вызовет панику
	// y := 10
	// w := reflect.ValueOf(y)
	// w.SetInt(20) // panic: reflect: reflect.Value.SetInt using unaddressable value
}

func getValueFromPointer() {
	fmt.Println("\n--- Работа с указателями ---")

	x := 100
	ptr := &x

	v := reflect.ValueOf(ptr)
	fmt.Printf("Тип указателя: %s\n", v.Type())
	fmt.Printf("Kind указателя: %s\n", v.Kind())

	// Получаем значение, на которое указывает указатель
	elem := v.Elem()
	fmt.Printf("Значение через Elem(): %v\n", elem.Interface())

	// Изменяем значение
	elem.SetInt(200)
	fmt.Printf("Значение x после изменения: %v\n", x)
}
