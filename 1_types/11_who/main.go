package main

import "fmt"

func main() {
	var x int = 42
	var y float64 = 3.14
	var s string = "hello"
	var err error = fmt.Errorf("oops")
	var p *int = nil

	fmt.Printf("x: %T\n", x)     // int
	fmt.Printf("y: %T\n", y)     // float64
	fmt.Printf("s: %T\n", s)     // string
	fmt.Printf("err: %T\n", err) // *errors.errorString (или другой тип, реализующий error)
	fmt.Printf("p: %T\n", p)     // *int

	//	Возможные ошибки
	var ptr *int = nil
	var i interface{} = ptr
	fmt.Printf("%T\n", i) // *int, а не nil
}
