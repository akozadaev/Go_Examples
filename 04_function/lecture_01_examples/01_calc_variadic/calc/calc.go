package calc

import "errors"

// ErrEmptyInput сообщает, что для вычисления нужно хотя бы одно значение.
var ErrEmptyInput = errors.New("empty input")

// Sum возвращает сумму всех чисел. Для пустого набора аргументов возвращает 0.
func Sum(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// Avg возвращает среднее арифметическое. Для пустого набора аргументов возвращает ErrEmptyInput.
func Avg(nums ...float64) (float64, error) {
	if len(nums) == 0 {
		return 0, ErrEmptyInput
	}

	total := 0.0
	for _, n := range nums {
		total += n
	}

	return total / float64(len(nums)), nil
}

// ApplyTwice применяет f к x два раза. Если f равна nil, возвращает x без изменений.
func ApplyTwice(f func(int) int, x int) int {
	if f == nil {
		return x
	}
	return f(f(x))
}
