package main

import "fmt"

func value(name string, result int) int {
	fmt.Println("value", name)
	return result
}

func main() {
	switch x := value("switch expr", 2); x {
	case value("case 1", 1):
		fmt.Println("one")
	case value("case 2a", 2), value("case 2b", 20):
		fmt.Println("two")
	default:
		fmt.Println("default")
	}
}
