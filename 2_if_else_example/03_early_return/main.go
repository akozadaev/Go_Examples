package main

import "fmt"

func enter(age int, hasPermission bool) {
	if age < 18 && !hasPermission {
		fmt.Println("access denied")
		return
	}

	fmt.Println("welcome")
}

func main() {
	enter(16, false)
	enter(20, false)
}
