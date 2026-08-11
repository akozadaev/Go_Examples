// Пакет memdemo содержит небольшие примеры для анализа утечек в кучу и измерения аллокаций.
package memdemo

import "bytes"

type Point struct {
	X int
	Y int
}

// SumLocal сохраняет данные локальными, если компилятор доказывает, что они не покидают функцию.
func SumLocal(x, y int) int {
	p := Point{X: x, Y: y}
	return p.X + p.Y
}

// NewPoint возвращает указатель. Место хранения объекта определяет escape analysis.
func NewPoint(x, y int) *Point {
	return &Point{X: x, Y: y}
}

// BuildMessageConcat намеренно создаёт промежуточные строки.
func BuildMessageConcat(parts []string) string {
	var result string
	for _, part := range parts {
		result += part
	}
	return result
}

// BuildMessageBuffer использует один растущий буфер.
func BuildMessageBuffer(parts []string) string {
	var buf bytes.Buffer
	for _, part := range parts {
		buf.WriteString(part)
	}
	return buf.String()
}
