package main

import "fmt"

func main() {
	fmt.Println(5 / 2)  // 2
	fmt.Println(-5 / 2) // -2

	var v int8 = 127
	v++ // v == -128
	fmt.Print(v, "\n")

	// Указатели
	x := 42
	p := &x // *int
	*p = 100
	fmt.Println(x) // 100

	p1 := new(int) // *p1 == 0
	p2 := new(42)  // *p2 == 42
	fmt.Println(p1, p2)
	fmt.Println(*p1, *p2)

	{
		type Person struct {
			Name string `json:"name"`
			Age  *int   `json:"age,omitempty"`
		}
		age := 30
		p := Person{Name: "Alice", Age: &age}
		fmt.Println(p)
		// нужно разыменовать
		fmt.Println(*p.Age)
	}

}

func stackVar() int {
	v := 42
	return v // на стеке
}
func heapVar() *int {
	v := 42
	return &v // утекает в кучу
}
