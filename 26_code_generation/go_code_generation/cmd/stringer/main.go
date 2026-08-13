package main

import (
	"fmt"

	"github.com/akozadaev/go-code-generation/stringer_example"
)

func main() {
	// Демонстрация использования stringer
	statuses := []stringer_example.Status{
		stringer_example.StatusPending,
		stringer_example.StatusInProgress,
		stringer_example.StatusCompleted,
		stringer_example.StatusCancelled,
	}

	fmt.Println("=== Демонстрация stringer ===")
	for _, status := range statuses {
		fmt.Printf("Status: %d, String: %s\n", status, status.String())
	}
}
