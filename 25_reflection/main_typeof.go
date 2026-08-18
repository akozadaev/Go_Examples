package main

import (
	"fmt"
	"reflect"
)

// main_typeof демонстрирует использование reflect.Type и TypeOf
func main() {
	fmt.Println("=== Примеры использования reflect.TypeOf ===")

	// 1. Получение типа переменной
	x := 42
	t := reflect.TypeOf(x)
	fmt.Printf("Тип переменной x: %s\n", t)

	// 2. Получение базового типа (Kind)
	fmt.Printf("Базовый тип (Kind) x: %s\n", t.Kind())

	// 3. Работа с различными типами
	demonstrateTypes()
}

func demonstrateTypes() {
	fmt.Println("\n--- Работа с различными типами ---")

	// int
	var i int = 10
	fmt.Printf("int: %s (Kind: %s)\n", reflect.TypeOf(i), reflect.TypeOf(i).Kind())

	// string
	var s string = "hello"
	fmt.Printf("string: %s (Kind: %s)\n", reflect.TypeOf(s), reflect.TypeOf(s).Kind())

	// bool
	var b bool = true
	fmt.Printf("bool: %s (Kind: %s)\n", reflect.TypeOf(b), reflect.TypeOf(b).Kind())

	// float64
	var f float64 = 3.14
	fmt.Printf("float64: %s (Kind: %s)\n", reflect.TypeOf(f), reflect.TypeOf(f).Kind())

	// slice
	var slice []int = []int{1, 2, 3}
	fmt.Printf("[]int: %s (Kind: %s)\n", reflect.TypeOf(slice), reflect.TypeOf(slice).Kind())

	// map
	var m map[string]int = map[string]int{"a": 1}
	fmt.Printf("map[string]int: %s (Kind: %s)\n", reflect.TypeOf(m), reflect.TypeOf(m).Kind())

	// struct
	type Person struct {
		Name string
		Age  int
	}
	var p Person = Person{"Alice", 30}
	fmt.Printf("Person: %s (Kind: %s)\n", reflect.TypeOf(p), reflect.TypeOf(p).Kind())

	// slice элементов структуры
	var people []Person = []Person{{"Bob", 25}}
	fmt.Printf("[]Person: %s (Kind: %s)\n", reflect.TypeOf(people), reflect.TypeOf(people).Kind())
}
