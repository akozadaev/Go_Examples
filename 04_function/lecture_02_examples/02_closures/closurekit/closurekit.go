package closurekit

// NewCounter возвращает функцию, которая увеличивает и возвращает собственный счетчик.
func NewCounter() func() int {
	count := 0
	return func() int {
		count++
		return count
	}
}

// MakeMultiplier возвращает функцию, которая умножает входное значение на factor.
func MakeMultiplier(factor int) func(int) int {
	return func(value int) int {
		return value * factor
	}
}

// Map применяет f к каждому значению и возвращает новый срез. Если f равна nil, значения копируются без изменений.
func Map(values []int, f func(int) int) []int {
	result := make([]int, len(values))
	for i, value := range values {
		if f == nil {
			result[i] = value
			continue
		}
		result[i] = f(value)
	}
	return result
}
