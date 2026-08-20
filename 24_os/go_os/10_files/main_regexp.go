package main

import (
	"fmt"
	"regexp"
)

func main() {
	re, err := regexp.Compile(`^[a-zA-Z0-9%+-]+@[a-zA-Z0-9]+\.[a-zA-Z]{2}$`)
	if err != nil {
		fmt.Println(fmt.Errorf("%v", err))
	}
	if re.MatchString("mail@mail.ru") {
		fmt.Println("match")
	} else {
		fmt.Println("not match")
	}

	re1, err := regexp.Compile(`([A-Z]+)=(\d+)`)
	if err != nil {
		fmt.Println(fmt.Errorf("%v", err))
	}
	math := re1.FindStringSubmatch("PORT=8080")
	if math != nil {
		for i, name := range re1.SubexpNames() {
			fmt.Println(i, name)

		}
	}
}
