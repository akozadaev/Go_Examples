package main

import "fmt"

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64 | ~string
}

func Sum[T Ordered](vals []T) T {
	var sum T
	for _, v := range vals {
		sum += v
	}
	return sum
}

func main() {
	// Целые числа
	ints := []int{1, 2, 3, 4, 5}
	fmt.Println("Sum of ints:", Sum(ints)) // 15

	// Числа с плавающей точкой
	floats := []float64{1.5, 2.5, 3.0}
	fmt.Println("Sum of floats:", Sum(floats)) // 7

	// Строки (конкатенация)
	strs := []string{"Hello", " ", "World", "!"}
	fmt.Println("Sum of strings:", Sum(strs)) // "Hello World!"

	// Пустой слайс
	empty := []int{}
	fmt.Println("Sum of empty:", Sum(empty)) // 0

	// Кастомный тип
	type MyInt int
	myInts := []MyInt{10, 20, 30}
	fmt.Println("Sum of MyInt:", Sum(myInts)) // 60
}
