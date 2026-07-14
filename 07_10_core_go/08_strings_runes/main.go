package main

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

func main() {
	fmt.Println("Тема 8. Строки и руны на задаче подготовки текста")
	fmt.Println()

	text := " Go, go! Привет, мир. 世界  "

	demonstrateBytesAndRunes(text)
	fmt.Println()

	words := extractWords(text)
	fmt.Println("2. Извлечение слов через unicode.IsLetter / unicode.IsDigit")
	fmt.Println("   слова:", words)
	fmt.Println()

	fmt.Println("3. Сборка отчета через strings.Builder")
	fmt.Println(buildWordLine(words))
	fmt.Println()

	demonstrateStandardStringFunctions(text)
	fmt.Println()

	demonstrateComparisonAndLiterals()
	fmt.Println()

	demonstrateConversionsAndSlicing()
	fmt.Println()

	demonstrateFormattingAndBytes()
	fmt.Println()

	demonstrateGraphemesAndMyths()
}

func demonstrateBytesAndRunes(text string) {
	fmt.Println("1. Строка хранит байты, range возвращает руны")
	fmt.Println("   исходный текст:", text)
	fmt.Println("   количество байт:", len(text))
	fmt.Println("   количество рун:", utf8.RuneCountInString(text))

	fmt.Print("   первые 12 байт через цикл по индексу:")
	for i := 0; i < len(text) && i < 12; i++ {
		fmt.Printf(" %02x", text[i])
	}
	fmt.Println()

	fmt.Println("   первые руны и их байтовые индексы:")
	shown := 0
	for byteIndex, r := range text {
		if unicode.IsSpace(r) {
			continue
		}
		fmt.Printf("     байт=%2d руна=%q\n", byteIndex, r)
		shown++
		if shown == 8 {
			break
		}
	}

	fmt.Printf("   первый байт строки: %d\n", text[0])
	fmt.Println("   вывод: индекс строки обращается к байту, а не к букве")
}

func extractWords(text string) []string {
	var words []string
	var b strings.Builder

	flush := func() {
		if b.Len() == 0 {
			return
		}
		words = append(words, strings.ToLower(b.String()))
		b.Reset()
	}

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			b.WriteRune(r)
			continue
		}
		flush()
	}
	flush()

	return words
}

func buildWordLine(words []string) string {
	var b strings.Builder
	b.Grow(len(words) * 12)

	b.WriteString("   нормализованная строка: ")
	for i, word := range words {
		if i > 0 {
			b.WriteString(" | ")
		}
		b.WriteString(word)
	}

	return b.String()
}

func demonstrateStandardStringFunctions(text string) {
	fmt.Println("4. Частые функции пакета strings")
	trimmed := strings.TrimSpace(text)
	fmt.Println("   TrimSpace:", trimmed)
	fmt.Println("   Contains(\"Привет\"):", strings.Contains(trimmed, "Привет"))
	fmt.Println("   HasPrefix(\"Go\"):", strings.HasPrefix(trimmed, "Go"))
	fmt.Println("   HasSuffix(\"世界\"):", strings.HasSuffix(trimmed, "世界"))
	fmt.Println("   Index(\"мир\"):", strings.Index(trimmed, "мир"))
	fmt.Println("   ReplaceAll запятых:", strings.ReplaceAll(trimmed, ",", ""))
	fmt.Println("   Fields:", strings.Fields(trimmed))
	fmt.Println("   Split по пробелу:", strings.Split(trimmed, " "))
	fmt.Println("   Join:", strings.Join([]string{"go", "строки", "руны"}, " + "))

	before, after, ok := strings.Cut(trimmed, "!")
	fmt.Printf("   Cut(\"!\"): до=%q после=%q найдено=%t\n", before, after, ok)
}

func demonstrateComparisonAndLiterals() {
	fmt.Println("5. Сравнение строк и строковые литералы")

	fmt.Println("   \"Go\" == \"go\":", "Go" == "go")
	fmt.Println("   EqualFold(\"Go\", \"go\"):", strings.EqualFold("Go", "go"))
	fmt.Println("   \"go\" < \"тест\":", "go" < "тест")
	fmt.Println("   strings.Compare(\"go\", \"тест\"):", strings.Compare("go", "тест"), "(обычно лучше использовать < или ==)")

	interpreted := "первая строка\nвторая строка"
	raw := `первая строка\nвторая строка`
	fmt.Printf("   интерпретируемый литерал: %q\n", interpreted)
	fmt.Printf("   сырой литерал: %q\n", raw)
}

func demonstrateConversionsAndSlicing() {
	fmt.Println("6. Преобразования и опасность срезов по байтам")

	s := "Привет"
	fmt.Printf("   %q: байт=%d рун=%d\n", s, len(s), utf8.RuneCountInString(s))
	fmt.Printf("   s[:1] дает байты % x, валидный UTF-8=%t\n", s[:1], utf8.ValidString(s[:1]))

	runes := []rune(s)
	fmt.Println("   первые три руны:", string(runes[:3]))

	mutable := []byte("статус=ok")
	mutable[len(mutable)-1] = 'K'
	fmt.Println("   []byte как изменяемый буфер:", string(mutable))
}

func demonstrateFormattingAndBytes() {
	fmt.Println("7. strconv, fmt и bytes")

	n, err := strconv.Atoi("42")
	if err != nil {
		fmt.Println("   ошибка strconv.Atoi:", err)
		return
	}
	fmt.Println("   strconv.Atoi(\"42\"):", n)
	fmt.Println("   strconv.Itoa(42):", strconv.Itoa(42))

	flag, err := strconv.ParseBool("true")
	if err != nil {
		fmt.Println("   ошибка strconv.ParseBool:", err)
		return
	}
	price, err := strconv.ParseFloat("19.95", 64)
	if err != nil {
		fmt.Println("   ошибка strconv.ParseFloat:", err)
		return
	}
	fmt.Println("   strconv.ParseBool(\"true\"):", flag)
	fmt.Println("   strconv.ParseFloat(\"19.95\", 64):", price)

	r := 'Ж'
	b := byte(0xff)
	fmt.Printf("   fmt: руна как символ=%c, руна как код=%U, байт как hex=%x\n", r, r, b)
	fmt.Println("   fmt.Sprintf:", fmt.Sprintf("%s: %d", "go", 3))
	fmt.Println("   fmt.Errorf:", fmt.Errorf("поле %q: %w", "name", strconv.ErrSyntax))

	payload := []byte("status=ok")
	fmt.Println("   bytes.Contains:", bytes.Contains(payload, []byte("ok")))
}

func demonstrateGraphemesAndMyths() {
	fmt.Println("8. Графемы и мифы про строки")

	family := "👨‍👩‍👧‍👦"
	combined := "е\u0301"
	precomposed := "é"

	fmt.Printf("   family emoji: байт=%d рун=%d видимый символ для человека обычно один\n", len(family), utf8.RuneCountInString(family))
	fmt.Printf("   combining form %q: байт=%d рун=%d\n", combined, len(combined), utf8.RuneCountInString(combined))
	fmt.Printf("   precomposed form %q: байт=%d рун=%d\n", precomposed, len(precomposed), utf8.RuneCountInString(precomposed))

	fmt.Println("   миф: len(s) считает символы -> нет, len(s) считает байты")
	fmt.Println("   миф: []rune всегда равно видимым символам -> нет, графема может состоять из нескольких рун")
	fmt.Println("   unsafe-преобразования []byte <-> string без копирования здесь не показываем: это опасная оптимизация для кода после профилирования")
}
