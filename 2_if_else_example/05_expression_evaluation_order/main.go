package main

import "fmt"

func value(name string, result int) int {
	fmt.Println("call", name)
	return result
}

func main() {
	x := value("A", 1) + value("B", 2) + value("C", 3)
	fmt.Println("x =", x)
}
