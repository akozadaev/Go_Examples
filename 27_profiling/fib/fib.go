// Пакет fib содержит намеренно различающиеся реализации для сравнения бенчмарками.
package fib

// Slow вычисляет число Фибоначчи рекурсивно и намеренно повторяет вычисления.
func Slow(n uint) uint64 {
	if n < 2 {
		return uint64(n)
	}
	return Slow(n-1) + Slow(n-2)
}

// Fast вычисляет число Фибоначчи итеративно за линейное время и с постоянной памятью.
func Fast(n uint) uint64 {
	var previous uint64
	current := uint64(1)
	for range n {
		previous, current = current, previous+current
	}
	return previous
}
