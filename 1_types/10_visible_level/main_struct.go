package main

import (
	"fmt"
	"mm/person"
)

func main() {
	p, err := person.NewPerson("Alexey", 45)
	fmt.Println(*p, err)
}
