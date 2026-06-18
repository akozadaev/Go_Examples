package main

import (
	"fmt"
	"math"
)

func main() {
	// Числа
	// var a int = 100
	var a int = 20000000000000000
	var b float64 = float64(a)
	fmt.Println("int -> float64:", b)
	fmt.Printf("int -> float64: %.0f\n", b)
	b = 20000000000000000
	var c int32 = int32(b) // Потенциально опасно
	fmt.Println("float64 -> int32:", c)

	// Байты и строки
	s := "Go"
	bts := []byte(s)
	fmt.Println("string -> []byte:", bts)
	fmt.Println("[]byte -> string:", string(bts))

	// Утверждение типа
	var i interface{} = 42
	if num, ok := i.(int); ok {
		fmt.Println("interface{} -> int:", num)
	}

	// Пользовательский тип
	type UserID int
	var id UserID = UserID(a)
	fmt.Println("int -> UserID:", id)

	val32to8, err := SafeInt32ToInt8(1000)
	if err != nil {
		fmt.Printf("Error: %v", err)
		fmt.Println()
	} else {
		fmt.Println("Value:", val32to8)
	}

	floatToInt64, err := SafeFloatToInt64(1000)
	if err != nil {
		fmt.Printf("Error: %v", err)
		fmt.Println()
	}
	fmt.Println("Value: ", floatToInt64)

}

func SafeInt32ToInt8(v int32) (int8, error) {
	if v > math.MaxInt8 || v < math.MinInt8 {
		return 0, fmt.Errorf("value %d out of int8 range", v)
	}
	return int8(v), nil
}

func SafeFloatToInt64(f float64) (int64, error) {
	if math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, fmt.Errorf("invalid float: %v", f)
	}
	if f > math.MaxInt64 || f < math.MinInt64 {
		return 0, fmt.Errorf("out of int64 range")
	}
	return int64(f), nil
}
