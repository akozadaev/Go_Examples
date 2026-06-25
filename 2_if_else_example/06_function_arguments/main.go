package main

import "fmt"

func arg(name string, value int) int {
	fmt.Println("arg", name)
	return value
}

func sum(a, b int) int {
	fmt.Println("inside sum")
	return a + b
}

func main() {
	result := sum(arg("first", 10), arg("second", 20))
	fmt.Println("result =", result)
}
