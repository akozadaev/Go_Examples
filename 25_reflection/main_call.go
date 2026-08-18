package main

import (
	"fmt"
	"reflect"
)

// Calculator представляет простой калькулятор
type Calculator struct {
	operations int
}

// Add выполняет сложение двух чисел
func (c *Calculator) Add(a, b int) int {
	c.operations++
	return a + b
}

// Subtract выполняет вычитание
func (c *Calculator) Subtract(a, b int) int {
	c.operations++
	return a - b
}

// Multiply выполняет умножение
func (c *Calculator) Multiply(a, b int) int {
	c.operations++
	return a * b
}

// SayHello - простая функция без получателя
func SayHello(name string) {
	fmt.Printf("Hello, %s!\n", name)
}

// Sum принимает произвольное количество аргументов
func Sum(numbers ...int) int {
	sum := 0
	for _, n := range numbers {
		sum += n
	}
	return sum
}

// main_call демонстрирует вызов функций и методов через рефлексию
func main() {
	fmt.Println("=== Примеры вызова функций через рефлексию ===")
	//
	//// 1. Вызов обычной функции
	//callFunction()
	//
	//// 2. Вызов метода структуры
	//callMethod()
	//
	//// 3. Вызов метода по имени
	//callMethodByName()

	// 4. Вызов функции с произвольным количеством аргументов
	callVariadicFunction()
}

func callFunction() {
	fmt.Println("\n--- Вызов обычной функции ---")

	funcValue := reflect.ValueOf(SayHello)
	fmt.Printf("Тип функции: %s\n", funcValue.Type())

	// Проверяем, что функция валидна
	if funcValue.IsValid() && funcValue.Kind() == reflect.Func {
		args := []reflect.Value{reflect.ValueOf("Alice")}
		funcValue.Call(args)
	}
}

func callMethod() {
	fmt.Println("\n--- Вызов метода структуры ---")

	calc := &Calculator{}
	i := calc.Add(5, 3) // Обычный вызов для увеличения счетчика
	fmt.Println(i)
	fmt.Printf("Операций выполнено: %d\n", calc.operations)

	// Вызов через рефлексию
	calcValue := reflect.ValueOf(calc)
	addMethod := calcValue.MethodByName("Add")

	if addMethod.IsValid() {
		args := []reflect.Value{reflect.ValueOf(10), reflect.ValueOf(20)}
		result := addMethod.Call(args)
		fmt.Printf("Результат вызова через рефлексию: %d\n", result[0].Int())
		fmt.Printf("Операций выполнено: %d\n", calc.operations)
	}
}

func callMethodByName() {
	fmt.Println("\n--- Вызов методов по имени ---")

	calc := &Calculator{}

	operations := []string{"Add", "Subtract", "Multiply"}
	values := [][]int{{10, 5}, {20, 7}, {3, 4}}

	calcValue := reflect.ValueOf(calc)

	for i, op := range operations {
		method := calcValue.MethodByName(op)
		if method.IsValid() {
			args := []reflect.Value{
				reflect.ValueOf(values[i][0]),
				reflect.ValueOf(values[i][1]),
			}
			result := method.Call(args)
			fmt.Printf("%s(%d, %d) = %d\n", op, values[i][0], values[i][1], result[0].Int())
		}
	}

	fmt.Printf("Всего операций выполнено: %d\n", calc.operations)
}

func callVariadicFunction() {
	fmt.Println("\n--- Вызов функции с произвольным количеством аргументов ---")

	sumFunc := reflect.ValueOf(Sum)
	fmt.Printf("Тип функции: %s\n", sumFunc.Type())

	if sumFunc.IsValid() && sumFunc.Kind() == reflect.Func {
		// Создаем слайс reflect.Value для аргументов
		args := []reflect.Value{
			reflect.ValueOf(1),
			reflect.ValueOf(2),
			reflect.ValueOf(3),
			reflect.ValueOf(4),
			reflect.ValueOf(5),
		}

		result := sumFunc.Call(args)
		fmt.Printf("Sum(1, 2, 3, 4, 5) = %d\n", result[0].Int())
	}

	// Вызов с CallSlice для variadic функций
	args := []reflect.Value{
		reflect.ValueOf([]int{10, 20, 30}),
	}
	result := sumFunc.CallSlice(args)
	fmt.Printf("Sum(10, 20, 30) = %d\n", result[0].Int())
}
