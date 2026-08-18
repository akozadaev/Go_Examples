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

type Student struct {
	Person
	University string
	Grade      float64
}

// main_struct демонстрирует работу со структурами через рефлексию
func main() {
	fmt.Println("=== Примеры работы со структурами ===")

	// 1. Перебор полей структуры
	iterateStructFields()

	// 2. Получение информации о полях
	getFieldInfo()

	// 3. Изменение значений полей
	modifyStructFields()

	// 4. Работа с вложенными структурами
	nestedStructs()
}

func iterateStructFields() {
	fmt.Println("\n--- Перебор полей структуры ---")

	p := Person{"Alice", 30, "alice@example.com", "New York"}
	v := reflect.ValueOf(p)
	t := v.Type()

	fmt.Printf("Структура: %s\n", t.Name())
	fmt.Printf("Количество полей: %d\n\n", t.NumField())

	// Go 1.26: Type.Fields и Value.Fields позволяют обходить поля итератором.
	fieldIndex := 0
	for fieldType, field := range v.Fields() {
		fmt.Printf("Поле %d:\n", fieldIndex)
		fmt.Printf("  Имя: %s\n", fieldType.Name)
		fmt.Printf("  Тип: %s\n", fieldType.Type)
		fmt.Printf("  Значение: %v\n", field.Interface())
		fmt.Println()
		fieldIndex++
	}
}

func getFieldInfo() {
	fmt.Println("--- Получение информации о полях ---")

	p := Person{"Bob", 25, "bob@example.com", "Los Angeles"}
	t := reflect.TypeOf(p)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fmt.Printf("Поле: %s\n", field.Name)
		fmt.Printf("  Тип: %s\n", field.Type)
		fmt.Printf("  PkgPath: %s\n", field.PkgPath)
		fmt.Println()
	}
}

func modifyStructFields() {
	fmt.Println("--- Изменение значений полей ---")

	p := Person{"Charlie", 35, "charlie@example.com", "Chicago"}
	v := reflect.ValueOf(&p).Elem()

	fmt.Printf("До изменения: %+v\n", p)

	// Изменяем поле Name
	nameField := v.FieldByName("Name")
	if nameField.IsValid() && nameField.CanSet() {
		nameField.SetString("Charlie Modified")
	}

	// Изменяем поле Age
	ageField := v.FieldByName("Age")
	if ageField.IsValid() && ageField.CanSet() {
		ageField.SetInt(36)
	}

	fmt.Printf("После изменения: %+v\n", p)
}

func nestedStructs() {
	fmt.Println("--- Работа с вложенными структурами ---")

	s := Student{
		Person: Person{
			Name:  "David",
			Age:   20,
			Email: "david@example.com",
			City:  "Boston",
		},
		University: "MIT",
		Grade:      3.8,
	}

	v := reflect.ValueOf(s)
	t := v.Type()

	fmt.Printf("Структура: %s\n", t.Name())
	fmt.Printf("Количество полей: %d\n\n", t.NumField())

	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		fieldValue := v.Field(i)

		fmt.Printf("Поле: %s\n", field.Name)
		fmt.Printf("  Тип: %s\n", field.Type)
		fmt.Printf("  Kind: %s\n", fieldValue.Kind())
		fmt.Printf("  Значение: %v\n", fieldValue.Interface())

		// Если это вложенная структура, показываем её поля
		if fieldValue.Kind() == reflect.Struct {
			fmt.Printf("  --- Вложенная структура ---\n")
			for j := 0; j < fieldValue.NumField(); j++ {
				nestedField := fieldValue.Type().Field(j)
				nestedValue := fieldValue.Field(j)
				fmt.Printf("    %s: %v\n", nestedField.Name, nestedValue.Interface())
			}
		}
		fmt.Println()
	}
}
