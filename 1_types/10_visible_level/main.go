package main

import (
	"fmt"
	"mm/math"
)

func main() {
	// Пакетный уровень
	fmt.Println(math.Pi)
	fmt.Println(math.Add(2, 4))
	//fmt.Println(math.internalVersion) // undefined: math.internalVersion
	//fmt.Println(math.multiply(3, 5)) // undefined: math.internalVersion

	// Пример локальной видимости
	x := 10
	if x > 0 {
		y := 5
		fmt.Println(x, y)
	}
	//fmt.Println(x, y) //undefined: y

}
