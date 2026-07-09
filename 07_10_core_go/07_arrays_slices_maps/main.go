package main

import (
	"fmt"
	"sort"
)

func main() {
	arrays()
	slices()
	maps()
}

func arrays() {
	numbers := [3]int{10, 20, 30}
	copyOfNumbers := numbers
	copyOfNumbers[0] = 99

	fmt.Println("исходный массив:", numbers)
	fmt.Println("копия массива:  ", copyOfNumbers)
}

func slices() {
	base := []string{"go", "тесты", "ошибки"}
	view := base[:2]
	view[0] = "Go"

	fmt.Println("базовый слайс после изменения представления:", base)
	fmt.Printf("представление: len=%d cap=%d значение=%v\n", len(view), cap(view), view)

	detached := append([]string(nil), view...)
	detached = append(detached, "юникод")
	detached[0] = "отдельный"

	fmt.Println("базовый слайс после изменения копии:", base)
	fmt.Println("отдельная копия слайса:", detached)
}

func maps() {
	words := []string{"go", "слайс", "go", "мапа", "слайс", "go"}
	counts := make(map[string]int)
	for _, word := range words {
		counts[word]++
	}

	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Println("стабильный вывод мапы:")
	for _, key := range keys {
		fmt.Printf("  %s -> %d\n", key, counts[key])
	}
}
