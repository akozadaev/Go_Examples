package main

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	text := "Go: Привет, 世界"

	fmt.Println("текст:", text)
	fmt.Println("байт:", len(text))
	fmt.Println("рун:", utf8.RuneCountInString(text))

	fmt.Println("range возвращает байтовый индекс и руну:")
	for i, r := range text {
		fmt.Printf("  байт=%2d руна=%q\n", i, r)
	}

	fmt.Println("только буквы:", keepLetters(text))
	fmt.Println("сборка через strings.Builder:", joinWithBuilder([]string{"массивы", "строки", "ошибки", "тесты"}))
}

func keepLetters(s string) string {
	var b strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) {
			b.WriteRune(unicode.ToLower(r))
		}
	}
	return b.String()
}

func joinWithBuilder(parts []string) string {
	var b strings.Builder
	for i, part := range parts {
		if i > 0 {
			b.WriteString(" -> ")
		}
		b.WriteString(part)
	}
	return b.String()
}
