package main

import (
	"fmt"
	"sort"
	"strings"
)

type WordStat struct {
	Word  string
	Count int
}

func main() {
	fmt.Println("Тема 7. Коллекции на задаче подсчета слов")
	fmt.Println()

	demonstrateArrays()
	fmt.Println()

	words := demonstrateSlices()
	fmt.Println()

	counts := countWords(words)
	demonstrateMaps(counts)
	fmt.Println()

	demonstrateRangeCopy()
}

func demonstrateArrays() {
	fmt.Println("1. Массив: фиксированный размер и копирование значения")

	topCounts := [3]int{5, 3, 2}
	copyOfTopCounts := topCounts
	copyOfTopCounts[0] = 99

	fmt.Println("   исходный массив: ", topCounts)
	fmt.Println("   измененная копия:", copyOfTopCounts)
	fmt.Println("   вывод: массив копируется целиком, исходное значение не изменилось")
}

func demonstrateSlices() []string {
	fmt.Println("2. Слайс: окно на общий массив, len/cap и безопасное копирование")

	words := []string{"go", "тесты", "ошибки", "go", "строки", "go"}
	firstThree := words[:3]

	fmt.Printf("   все слова: %v\n", words)
	fmt.Printf("   первые три: %v, len=%d, cap=%d\n", firstThree, len(firstThree), cap(firstThree))

	firstThree[0] = "Go"
	fmt.Println("   после изменения firstThree[0] изменился и исходный слайс:", words)

	limited := words[:3:3]
	limited = append(limited, "юникод")
	fmt.Printf("   append к words[:3:3]: %v, len=%d, cap=%d\n", limited, len(limited), cap(limited))
	fmt.Println("   исходный слайс после безопасного append:", words)

	detached := append([]string(nil), words[:3]...)
	detached[0] = "копия"
	fmt.Println("   независимая копия:", detached)
	fmt.Println("   исходный слайс после изменения копии:", words)

	return words
}

func countWords(words []string) map[string]int {
	counts := make(map[string]int, len(words))
	for _, word := range words {
		counts[strings.ToLower(word)]++
	}
	return counts
}

func demonstrateMaps(counts map[string]int) {
	fmt.Println("3. Мапа: подсчет частот и стабильный вывод через сортировку ключей")

	fmt.Println("   чтение отсутствующего ключа дает нулевое значение:", counts["нет"])

	value, ok := counts["go"]
	fmt.Printf("   проверка наличия ключа \"go\": значение=%d, есть=%t\n", value, ok)

	delete(counts, "ошибки")
	fmt.Println("   после delete(\"ошибки\") ключ можно безопасно читать:", counts["ошибки"])

	keys := make([]string, 0, len(counts))
	for key := range counts {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	fmt.Println("   стабильный отчет:")
	for _, key := range keys {
		fmt.Printf("     %s: %d\n", key, counts[key])
	}
}

func demonstrateRangeCopy() {
	fmt.Println("4. range по слайсу: значение в цикле является копией")

	stats := []WordStat{
		{Word: "go", Count: 3},
		{Word: "тесты", Count: 1},
	}

	for _, stat := range stats {
		stat.Count = 100
	}
	fmt.Println("   после изменения копии в range:", stats)

	for i := range stats {
		stats[i].Count = 100
	}
	fmt.Println("   после изменения по индексу:", stats)
}
