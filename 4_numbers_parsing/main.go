package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	f, err := strconv.ParseFloat("NaN", 64)
	fmt.Print(f, err, "\n")
	// err == nil, но f - это NaN!

	vals := []string{"0.15", "1.14", "1.15", "0.29"}

	for _, v := range vals {
		f, _ := strconv.ParseFloat(v, 64)
		scaled := f * 100 // Ваше умножение

		// float64 печатается с округлением для читаемости,
		// но внутри хранится неточное значение!
		fmt.Printf("String: %-5s | float64 * 100: %-22v | int64: %d\n",
			v, scaled, int64(scaled))
	}

	val64, err := ParseAndScaleFloatNaive("1.14", 100)
	fmt.Print("Ошибка: 1.14 * 100 != ", val64, err, "\n")
	val64, err = ParseAndScaleFloat("1.14", 2)
	fmt.Print(val64, err, "\n")
}

// Наивная реализация (Опасна!)
func ParseAndScaleFloatNaive(s string, scale int) (int64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, err
	}

	// ПРОБЛЕМА: float64 не может точно представить 0.1, 0.2, 1.15 и т.д.
	// 1.15 * 100 может дать 114.99999999999999
	// scaled := f * math.Pow10(scale)
	scaled := f * float64(scale)
	fmt.Println(scaled)

	// При приведении к int64 дробная часть просто отбросится!
	// int64(114.99999999999999) == 114 (ОШИБКА!)
	return int64(scaled), nil
}

// С математическим округлением
func ParseAndScaleFloat(s string, scale int) (int64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("parse float: %w", err)
	}

	// 1. Отсекаем недопустимые значения
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, fmt.Errorf("invalid float value: %v", f)
	}

	// 2. Вычисляем множитель
	multiplier := math.Pow10(scale)

	// 3. Проверка на переполнение ДО умножения (опционально, но полезно)
	if f > math.MaxInt64/multiplier || f < math.MinInt64/multiplier {
		return 0, fmt.Errorf("value %f is out of int64 range after scaling", f)
	}

	// 4. Умножаем и ОКРУГЛЯЕМ до ближайшего целого
	scaled := math.Round(f * multiplier)

	return int64(scaled), nil
}
