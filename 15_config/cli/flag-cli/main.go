package main

import (
	"flag"
	"fmt"
)

func main() {
	// Объявляем переменные для хранения значений флагов
	var name string
	var age int
	var verbose bool

	// Регистрируем флаги
	flag.StringVar(&name, "name", "World", "Имя для приветствия")
	flag.StringVar(&name, "n", "World", "Имя для приветствия (короткая форма)")

	flag.IntVar(&age, "age", 0, "Возраст пользователя")
	flag.IntVar(&age, "a", 0, "Возраст пользователя (короткая форма)")

	flag.BoolVar(&verbose, "verbose", false, "Включить подробный вывод")
	flag.BoolVar(&verbose, "v", false, "Включить подробный вывод (короткая форма)")

	// Парсим флаги из командной строки
	flag.Parse()

	// Выводим информацию в зависимости от флагов
	if verbose {
		fmt.Println("Режим отладки включён")
		fmt.Printf("Получено имя: %s\n", name)
		fmt.Printf("Получен возраст: %d\n", age)
	}

	// Основная логика
	if age > 0 {
		fmt.Printf("Привет, %s! Тебе %d лет.\n", name, age)
	} else {
		fmt.Printf("Привет, %s!\n", name)
	}

	// Вывод оставшихся аргументов (если есть)
	args := flag.Args()
	if len(args) > 0 {
		fmt.Printf("Дополнительные аргументы: %v\n", args)
	}
}
