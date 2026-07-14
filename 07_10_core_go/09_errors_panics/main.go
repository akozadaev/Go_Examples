package main

import (
	"errors"
	"fmt"
	"strings"
)

var ErrEmptyText = errors.New("текст пуст")

type ValidationError struct {
	Field string
	Value string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("некорректное поле %s: %q", e.Field, e.Value)
}

func main() {
	fmt.Println("Тема 9. Ошибки на задаче подготовки отчета")
	fmt.Println()

	for _, text := range []string{
		"Go, go, тесты!",
		"   ",
		"Go",
	} {
		report, err := buildReport(text)
		if err != nil {
			fmt.Printf("вход %q: %v\n", text, err)

			if errors.Is(err, ErrEmptyText) {
				fmt.Println("  причина найдена через errors.Is: пустой текст")
			}

			var validationErr ValidationError
			if errors.As(err, &validationErr) {
				fmt.Printf("  детали найдены через errors.As: поле=%s значение=%q\n", validationErr.Field, validationErr.Value)
			}
			fmt.Println()
			continue
		}

		fmt.Println(report)
		fmt.Println()
	}

	result, err := safelyDivide(10, 0)
	fmt.Println("recover на границе опасной операции:", result, err)
}

func buildReport(text string) (report string, err error) {
	defer fmt.Println("  defer: завершена попытка построить отчет")

	words, err := parseWords(text)
	if err != nil {
		return "", fmt.Errorf("подготовка слов: %w", err)
	}

	counts := make(map[string]int, len(words))
	for _, word := range words {
		counts[word]++
	}

	return fmt.Sprintf("отчет для %q: слов=%d уникальных=%d", text, len(words), len(counts)), nil
}

func parseWords(text string) ([]string, error) {
	if strings.TrimSpace(text) == "" {
		return nil, ErrEmptyText
	}

	words := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return r == ' ' || r == ',' || r == '!' || r == '.'
	})

	if len(words) < 2 {
		return nil, ValidationError{
			Field: "количество слов",
			Value: fmt.Sprint(len(words)),
		}
	}

	return words, nil
}

func safelyDivide(a, b int) (result int, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("паника перехвачена: %v", r)
		}
	}()

	if b == 0 {
		panic("деление на ноль")
	}

	return a / b, nil
}
