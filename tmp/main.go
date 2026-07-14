package main

import (
	"fmt"
	"strings"
)

func main() {
	text := "Go, go! Привет"
	firstByte := text[0]
	fmt.Println(firstByte)

	//	==================

	for i := 0; i < len(text); i++ {
		fmt.Printf("%02x ", text[i])
	}

	//	==================
	fmt.Println()
	for byteIndex, r := range text {
		fmt.Printf("Индекс байта: %2d | Символ (руна): %c | Код Unicode: %U | Тип: %T\n", byteIndex, r, r, r)
	}

	// Количество рун:
	en := "r"
	ru := "й"

	wave := "👋"         // 1 руна
	waveSkin := "👋🏽"    // 2 руны: машущая рука + модификатор кожи
	family := "👨‍👩‍👧‍👦" // 7 рун, соединённых ZWJ (Zero Width Joiner)

	fmt.Printf("%s - рун: %d, байт: %d\n", en, len([]rune(en)), len(en))
	fmt.Printf("%s - рун: %d, байт: %d\n", ru, len([]rune(ru)), len(ru))
	fmt.Printf("%s - рун: %d, байт: %d\n", wave, len([]rune(wave)), len(wave))
	fmt.Printf("%s - рун: %d, байт: %d\n", waveSkin, len([]rune(waveSkin)), len(waveSkin))
	fmt.Printf("%s - рун: %d, байт: %d\n", family, len([]rune(family)), len(family))

	//	==================

	// Строковый литерал
	line := "первая строка\nвторая строка"
	fmt.Println(line)
	line = "первая строка\tвторая строка после табуляции"
	fmt.Println(line)

	path := "d:\tmp\report.txt"
	fmt.Println("===path===")
	fmt.Println(path)

	path = "d:\\tmp\\report.txt"
	fmt.Println("===path2===")
	fmt.Println(path)

	path = `d:\tmp\report.txt`
	fmt.Println("===path3===")
	fmt.Println(path)

	sql := `
SELECT avatar_id
FROM profiles
WHERE id = 1
  AND profiles.deleted_at IS NULL`
	fmt.Println("===sql===")
	fmt.Println(sql)

	//	[]byte и []rune
	runes := []rune("Привет")
	runes[0] = 'п'
	fmt.Println(string(runes))

	//	=============
	bytes := []byte(text)
	sameText := string(bytes)

	fmt.Println(sameText)

	runes = []rune(text)
	textAgain := string(runes)

	fmt.Println(textAgain)

	//	Срезы строк
	s := "привет"
	part := s[:2]
	fmt.Println(part)

	// Сравнение строк
	a := "стр1"
	b := "стр2"

	if a < b {
		fmt.Println("a идет раньше b")
	} else if a == b {
		fmt.Println("строки равны")
	} else {
		fmt.Println("a идет позже b")
	}
	// CСравнение без учета регистра
	fmt.Println(strings.EqualFold("Go", "go"))
}
