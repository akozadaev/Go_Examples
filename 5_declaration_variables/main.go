package main

import "fmt"

func main() {
	{
		var global int
		var x int = 10
		var y = 20   // вывод типа
		z := 30      // только внутри функций
		a, b := 1, 2 // множественное
		fmt.Print(global, x, y, z, a, b, "\n")
	}

	{
		x, y := 1, 2
		x, z := 3, 4 // OK (z новая)
		fmt.Print(x, y, z, "\n")
		x, y = y, x // обмен без временной переменной
		fmt.Print(x, y, z, "\n")
	}

	{
		const Pi = 3.14
		const (
			A = iota // 0
			B        // 1
			C        // 2
		)
		fmt.Print(Pi, A, B, C, "\n")
	}

	{
		
	}
}
