package main

import "fmt"

// Интерфейс описывает контракт - любое поведение, которое умеет Greet()
type Greeter interface {
	Greet() string
}

// Реализация 1: EnglishGreeter
type EnglishGreeter struct {
	Name string
}

func (g EnglishGreeter) Greet() string {
	return "Hello, " + g.Name + "!"
}

// Реализация 2: RussianGreeter
type RussianGreeter struct {
	Name string
}

func (g RussianGreeter) Greet() string {
	return "Привет, " + g.Name + "!"
}

// Функция принимает интерфейс - работает с любой реализацией
func printGreeting(g Greeter) {
	fmt.Println(g.Greet())
}

func main() {
	// Полиморфизм: одна функция, разное поведение
	english := EnglishGreeter{Name: "Alice"}
	russian := RussianGreeter{Name: "Борис"}

	printGreeting(english) // Hello, Alice!
	printGreeting(russian) // Привет, Борис!

	// Интерфейс можно присвоить переменной
	var greeter Greeter
	greeter = english
	fmt.Println(greeter.Greet()) // Hello, Alice!

	greeter = russian
	fmt.Println(greeter.Greet()) // Привет, Борис!
}
