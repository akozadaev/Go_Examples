package main

import (
	"fmt"
	"github.com/spf13/pflag"
)

func main() {
	var name string
	var age int

	pflag.StringVarP(&name, "name", "n", "Guest", "Your name")
	pflag.IntVarP(&age, "age", "a", 30, "Your age")

	pflag.Parse()

	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Age: %d\n", age)
}
