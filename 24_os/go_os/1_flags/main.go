package main

import (
	"flag"
	"fmt"
)

func main() {
	var name string
	var age int

	flag.StringVar(&name, "name", "Guest", "Your name")
	flag.IntVar(&age, "age", 30, "Your age")

	flag.Parse()
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Age: %d\n", age)

}

func parseArg(args []string) (string, error) {
	fs := flag.NewFlagSet("app", flag.ContinueOnError)
	name := fs.String("name", "", "Your name")
	err := fs.Parse(args)
	if err != nil {
		return "", err
	}
	return *name, err
}
