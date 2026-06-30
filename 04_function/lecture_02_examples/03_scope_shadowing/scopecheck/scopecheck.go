package scopecheck

import (
	"errors"
	"strconv"
	"strings"
)

// ErrEmptyToken сообщает о пустом элементе в списке целых чисел, разделенных запятыми.
var ErrEmptyToken = errors.New("empty token")

// CountPositiveWithBug демонстрирует затенение: внутренний count не является внешним count.
func CountPositiveWithBug(values []int) int {
	count := 0
	for _, value := range values {
		if value > 0 {
			count := count + 1
			_ = count
		}
	}
	return count
}

// CountPositive считает положительные значения без затенения count.
func CountPositive(values []int) int {
	count := 0
	for _, value := range values {
		if value > 0 {
			count = count + 1
		}
	}
	return count
}

// ParseCSVInts разбирает целые числа, разделенные запятыми, и возвращает новый срез.
func ParseCSVInts(input string) ([]int, error) {
	if strings.TrimSpace(input) == "" {
		return []int{}, nil
	}

	parts := strings.Split(input, ",")
	values := make([]int, 0, len(parts))

	for _, part := range parts {
		token := strings.TrimSpace(part)
		if token == "" {
			return nil, ErrEmptyToken
		}

		value, err := strconv.Atoi(token)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}

	return values, nil
}

// FindLabel использует переменные из области видимости if, не выпуская их за пределы блока.
func FindLabel(labels map[string]string, key string) string {
	if label, ok := labels[key]; ok {
		return label
	}
	return "unknown"
}
