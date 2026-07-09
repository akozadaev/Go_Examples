package main

import (
	"errors"
	"fmt"
)

var ErrUserNotFound = errors.New("пользователь не найден")

type ValidationError struct {
	Field string
	Value string
}

func (e ValidationError) Error() string {
	return fmt.Sprintf("некорректное поле %s: %q", e.Field, e.Value)
}

func main() {
	for _, id := range []int{42, -1, 7} {
		name, err := loadUserName(id)
		if err != nil {
			fmt.Printf("загрузка пользователя %d: %v\n", id, err)

			if errors.Is(err, ErrUserNotFound) {
				fmt.Println("  найдена сигнальная ошибка")
			}

			var validationErr ValidationError
			if errors.As(err, &validationErr) {
				fmt.Printf("  найдена ошибка валидации: поле=%s\n", validationErr.Field)
			}
			continue
		}

		fmt.Println("загружен:", name)
	}

	result, err := safeDivide(10, 0)
	fmt.Println("безопасное деление:", result, err)
}

func loadUserName(id int) (string, error) {
	defer fmt.Println("  очистка после загрузки пользователя с id", id)

	if id < 0 {
		return "", fmt.Errorf("проверка id: %w", ValidationError{
			Field: "идентификатор",
			Value: fmt.Sprint(id),
		})
	}

	if id != 42 {
		return "", fmt.Errorf("запрос к хранилищу: %w", ErrUserNotFound)
	}

	return "Ада", nil
}

func safeDivide(a, b int) (result int, err error) {
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
