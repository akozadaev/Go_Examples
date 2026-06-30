package transform

import "strings"

// Step преобразует строку и возвращает преобразованное значение.
type Step func(string) string

// Pipeline применяет шаги преобразования слева направо.
func Pipeline(input string, steps ...Step) string {
	result := input
	for _, step := range steps {
		if step == nil {
			continue
		}
		result = step(result)
	}
	return result
}

// Trim удаляет пробельные символы в начале и конце строки.
func Trim(s string) string {
	return strings.TrimSpace(s)
}

// Lower переводит строку в нижний регистр.
func Lower(s string) string {
	return strings.ToLower(s)
}

// Prefix возвращает Step, который добавляет prefix в начало строки.
func Prefix(prefix string) Step {
	return func(s string) string {
		return prefix + s
	}
}
