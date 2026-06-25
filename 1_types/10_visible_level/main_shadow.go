package main

import "fmt"

var level = "global"

func main() {
	fmt.Println(level)
	level = "local"
	fmt.Println(level)

	if true {
		level := "block"
		fmt.Println(level)
	}
	
	fmt.Println(level)
}
