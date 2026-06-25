package main

import "fmt"

func printValue(v any) {
	switch value := v.(type) {
	case nil:
		fmt.Println("nil")
	case int:
		fmt.Println("int:", value)
	case string:
		fmt.Println("string:", value)
	case bool:
		fmt.Println("bool:", value)
	default:
		fmt.Printf("unknown type %T\n", value)
	}
}

func main() {
	printValue(10)
	printValue("hello")
	printValue(nil)
}
