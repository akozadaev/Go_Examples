package main

import (
	"errors"
	"fmt"
	"strings"
)

func main() {
	/*a := 10
	b := 15
	total := a + b
	fmt.Println(total)

	fmt.Println(Add(a, b))
	q, r := DivMod(10, 3)
	fmt.Println(q, r)

	d, err := Divide(4, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(d)

	fmt.Println(Avg())
	fmt.Println(Avg(1, 2, 3))

	fmt.Println(Avg1())
	fmt.Println(Avg1(1, 2, 3))
	BuildSummary(1)
	values := []float64{1, 2, 3}
	BuildSummary(values...)

	double := func(x int) int {
		return x * 2
	}
	result := ApplyTwice(double, 2)
	fmt.Println(result)
	*/

	/*	result := Pipeline(" LESSON  ", strings.TrimSpace, strings.ToLower)
		fmt.Println(result)*/

	/*	normalize := func(s string) string {
			return strings.ToLower(strings.TrimSpace(s))
		}
		fmt.Println(normalize(" LESSON  "))

	*/
	var debug bool
	debug = true

	message := func(d bool) string {
		if debug {
			return "debug mode"
		}
		return "normal mode"
	}(debug)

	fmt.Println(message)
}

func Add(a int, b int) int {
	return a + b
}

func PrepareUserName(input string) string {
	trimmed := strings.TrimSpace(input)
	return strings.ToLower(trimmed)
}
func NormalizeUserName(input string) string {
	return strings.ToLower(strings.TrimSpace(input))
}

func Square(x int) int {
	return x * x
}

func DivMod(a, b int) (int, int) {
	return a / b, a % b
}

func Divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("division by zero")
	}
	return a / b, nil
}

func Avg(nums ...float64) float64 {
	total := 0.0
	for _, n := range nums {
		total += n
	}
	return total / float64(len(nums))
}

var ErrEmptyInput = errors.New("no values")

func Avg1(nums ...float64) (float64, error) {
	if len(nums) == 0 {
		return 0, ErrEmptyInput
	}
	total := 0.0
	for _, n := range nums {
		total += n
	}
	return total / float64(len(nums)), nil
}

func BuildSummary(scores ...float64) (min float64, max float64, avg float64, ok bool) {
	fmt.Println(fmt.Sprintf("%T", scores))
	if len(scores) == 0 {
		return 0, 0, 0, false
	}
	min = scores[0]
	max = scores[0]
	total := 0.0
	for _, score := range scores {
		if score < min {
			min = score
		}
		if score > max {
			max = score
		}
		total += score

	}
	avg = total / float64(len(scores))
	ok = true
	return //min, max, avg, ok
}

func ApplyTwice(f func(int) int, x int) int {
	return f(f(x))
}

type Step func(string) string

func Pipeline(input string, steps ...Step) string {
	result := input
	for _, step := range steps {
		result = step(result)

	}
	return result
}
