package main

import (
	"fmt"

	"go_examples/23_memory/memdemo"
)

func main() {
	fmt.Println("local sum:", memdemo.SumLocal(20, 22))
	fmt.Println("returned point:", memdemo.NewPoint(20, 22))
}
