package main

import "fmt"

func statusByScore(score int) string {
	if message := fmt.Sprintf("score = %d", score); score >= 60 {
		fmt.Println(message)
		return "passed"
	} else {
		fmt.Println(message)
		return "failed"
	}
}

func main() {
	fmt.Println(statusByScore(75))
	fmt.Println(statusByScore(40))
}
