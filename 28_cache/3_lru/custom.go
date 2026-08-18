//go:build experiment

package main

import (
	"fmt"
)

type casche struct {
	key   string
	val   int
	order int
	cnt   int
}

func main() {
	var c [3]casche
	fmt.Println("Исходный массив:", c)

	for i := range c {
		c[i] = casche{
			key:   string('a' + i),
			val:   i * 10, // произвольное значение
			order: i + 1,  // номер позиции (начиная с 1)
			cnt:   len(c), // длина массива
		}
		fmt.Println("Заполняем элемент:", c[i])
	}

	fmt.Println("Заполненный массив:", c)
}
