package main

import (
	"fmt"
	"math"
	"os"
	"strconv"
	"text/tabwriter"
)

func main() {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)

	fmt.Fprintln(w, "Тип\tРазмер (бит)\tДиапазон")
	fmt.Fprintln(w, "---\t-------------\t--------")

	// Знаковые целые
	fmt.Fprintf(w, "int8\t8\t%d … %d\n", int8(math.MinInt8), int8(math.MaxInt8))
	fmt.Fprintf(w, "int16\t16\t%d … %d\n", int16(math.MinInt16), int16(math.MaxInt16))
	fmt.Fprintf(w, "int32\t32\t%d … %d\n", int32(math.MinInt32), int32(math.MaxInt32))
	fmt.Fprintf(w, "int64\t64\t%d … %d\n", int64(math.MinInt64), int64(math.MaxInt64))
	fmt.Fprintf(w, "int\t%d\t%d … %d\n", strconv.IntSize, int(math.MinInt), int(math.MaxInt))

	// Беззнаковые целые - ВАЖНО: явное приведение к uint*
	fmt.Fprintf(w, "uint8\t8\t%d … %d\n", uint8(0), uint8(math.MaxUint8))
	fmt.Fprintf(w, "uint16\t16\t%d … %d\n", uint16(0), uint16(math.MaxUint16))
	fmt.Fprintf(w, "uint32\t32\t%d … %d\n", uint32(0), uint32(math.MaxUint32))
	fmt.Fprintf(w, "uint64\t64\t%d … %d\n", uint64(0), uint64(math.MaxUint64)) // ← исправлено
	fmt.Fprintf(w, "uint\t%d\t%d … %d\n", strconv.IntSize, uint(0), uint(math.MaxUint))

	// Числа с плавающей точкой
	fmt.Fprintf(w, "float32\t32\t±%e … ±%e\n", math.SmallestNonzeroFloat32, math.MaxFloat32)
	fmt.Fprintf(w, "float64\t64\t±%e … ±%e\n", math.SmallestNonzeroFloat64, math.MaxFloat64)

	// Комплексные и алиасы
	fmt.Fprintln(w, "complex64\t64\tсостоит из двух float32")
	fmt.Fprintln(w, "complex128\t128\tсостоит из двух float64")
	fmt.Fprintln(w, "byte\t8\t= uint8 (0 … 255)")
	fmt.Fprintln(w, "rune\t32\t= int32 (Unicode code point)")

	w.Flush()
}
