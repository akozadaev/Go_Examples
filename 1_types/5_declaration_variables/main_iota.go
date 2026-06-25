package main

import "fmt"

type Status int

const (
	StatusUnknown Status = iota
	StatusPending
	StatusActive
	StatusClosed
)

// Добавим String() метод для красивого вывода
func (s Status) String() string {
	switch s {
	case StatusPending:
		return "PENDING"
	case StatusActive:
		return "ACTIVE"
	case StatusClosed:
		return "CLOSED"
	default:
		return "UNKNOWN"
	}
}

func main() {
	s := StatusActive
	fmt.Println(s)      // "ACTIVE"
	fmt.Println(int(s)) // 2
}
