package main

import "fmt"

type MathFunc func(a, b int) int

func apply(f MathFunc, a, b int) int {
	return f(a, b)
}

func main() {
	add := func(a, b int) int {
		return a + b
	}
	result := apply(add, 1, 1)
	fmt.Println(result)
}
