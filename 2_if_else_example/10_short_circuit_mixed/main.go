package main

import "fmt"

func f(name string, result bool) bool {
	fmt.Println(name)
	return result
}

func main() {
	if f("A", false) && f("B", true) || f("C", true) {
		fmt.Println("ok")
	}
}
