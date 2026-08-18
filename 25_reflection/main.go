package main

import (
	"fmt"
	"reflect"
)

type Person struct {
	Name string
	Age  int
}

// main демонстрирует базовые возможности пакета reflect
func main() {
	fmt.Println("=== Базовые примеры reflect ===")

	// 1. reflect.Type и TypeOf
	demonstrateType()

	// 2. reflect.Value и ValueOf
	demonstrateValue()

	// 3. Изменение значения
	demonstrateModify()

	// 4. Работа со структурами
	demonstrateStruct()
}

func demonstrateType() {
	fmt.Println("\n--- reflect.Type и TypeOf ---")

	var i int = 42
	t := reflect.TypeOf(i)
	fmt.Printf("Тип переменной: %s\n", t)
	fmt.Printf("Kind: %s\n", t.Kind())

	var s string = "hello"
	ts := reflect.TypeOf(s)
	fmt.Printf("Тип переменной: %s\n", ts)
	fmt.Printf("Kind: %s\n", ts.Kind())
}

func demonstrateValue() {
	fmt.Println("\n--- reflect.Value и ValueOf ---")

	var i int = 42
	v := reflect.ValueOf(i)
	fmt.Printf("Значение: %v\n", v.Interface())
	fmt.Printf("Int значение: %d\n", v.Int())

	var s string = "hello"
	vs := reflect.ValueOf(s)
	fmt.Printf("Значение: %v\n", vs.Interface())
	fmt.Printf("String значение: %s\n", vs.String())
}

func demonstrateModify() {
	fmt.Println("\n--- Изменение значения ---")

	x := 42
	fmt.Printf("До изменения: %d\n", x)

	v := reflect.ValueOf(&x) // Важно: получить указатель
	v.Elem().SetInt(24)

	fmt.Printf("После изменения: %d\n", x)
}

func demonstrateStruct() {
	fmt.Println("\n--- Работа со структурами ---")

	p := Person{"Alice", 30}
	v := reflect.ValueOf(p)
	t := v.Type()

	fmt.Printf("Структура: %s\n", t.Name())
	fmt.Printf("Количество полей: %d\n", t.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		fmt.Printf("%s: %v\n", fieldType.Name, field.Interface())
	}
}
