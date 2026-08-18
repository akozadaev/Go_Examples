package main

import (
	"fmt"
	"unsafe"
)

type ExampleStruct struct {
	A int32
	B int64
	C int32
}

// main_unsafe демонстрирует использование пакета unsafe
func main() {
	fmt.Println("=== Примеры использования unsafe ===")

	// 1. Получение размера и выравнивания
	getSizeAndAlignment()

	// 2. Допустимые преобразования указателей
	convertPointers()

	// 3. Операции с памятью
	memoryOperations()

	// 4. Структуры и память
	structMemory()
}

func getSizeAndAlignment() {
	fmt.Println("\n--- Размер и выравнивание типов ---")

	var i int32
	var i64 int64
	var s string
	var slice []int
	var m map[string]int

	fmt.Printf("int32: размер = %d байт, выравнивание = %d байт\n",
		unsafe.Sizeof(i), unsafe.Alignof(i))
	fmt.Printf("int64: размер = %d байт, выравнивание = %d байт\n",
		unsafe.Sizeof(i64), unsafe.Alignof(i64))
	fmt.Printf("string: размер = %d байт, выравнивание = %d байт\n",
		unsafe.Sizeof(s), unsafe.Alignof(s))
	fmt.Printf("[]int: размер = %d байт, выравнивание = %d байт\n",
		unsafe.Sizeof(slice), unsafe.Alignof(slice))
	fmt.Printf("map[string]int: размер = %d байт, выравнивание = %d байт\n",
		unsafe.Sizeof(m), unsafe.Alignof(m))
}

func convertPointers() {
	fmt.Println("\n--- Преобразование указателей ---")

	var x int32 = 42
	var y float64 = 3.14

	// Получаем указатели
	xPtr := unsafe.Pointer(&x)
	yPtr := unsafe.Pointer(&y)

	fmt.Printf("Значение x: %d\n", x)
	fmt.Printf("Значение y: %f\n", y)

	// Читаем значение через unsafe.Pointer
	xValue := *(*int32)(xPtr)
	fmt.Printf("Значение x через unsafe.Pointer: %d\n", xValue)

	// Правильное преобразование для y
	yValue := *(*float64)(yPtr)
	fmt.Printf("Значение y через unsafe.Pointer: %f\n", yValue)

	// Нельзя разыменовывать xPtr как *float64: float64 больше int32,
	// поэтому такое чтение вышло бы за границы объекта x. unsafe отключает
	// проверку типов, но не делает недопустимый доступ к памяти корректным.
	fmt.Println("Чтение xPtr как *float64 намеренно не выполняется: это небезопасно")
}

func memoryOperations() {
	fmt.Println("\n--- Операции с памятью ---")

	// Offsetof показывает смещение поля в структуре
	type Point struct {
		X int32
		Y int32
		Z int32
	}

	p := Point{X: 10, Y: 20, Z: 30}

	fmt.Printf("Размер структуры Point: %d байт\n", unsafe.Sizeof(p))
	fmt.Printf("Смещение поля X: %d байт\n", unsafe.Offsetof(p.X))
	fmt.Printf("Смещение поля Y: %d байт\n", unsafe.Offsetof(p.Y))
	fmt.Printf("Смещение поля Z: %d байт\n", unsafe.Offsetof(p.Z))
}

func structMemory() {
	fmt.Println("\n--- Структуры и память ---")

	s := ExampleStruct{
		A: 100,
		B: 200,
		C: 300,
	}

	fmt.Printf("Размер ExampleStruct: %d байт\n", unsafe.Sizeof(s))
	fmt.Printf("Выравнивание: %d байт\n", unsafe.Alignof(s))

	// Получаем указатель на структуру
	structPtr := unsafe.Pointer(&s)

	// Получаем указатель на поле A
	aOffset := unsafe.Offsetof(s.A)
	aPtr := unsafe.Add(structPtr, aOffset)
	aValue := *(*int32)(aPtr)
	fmt.Printf("Поле A через указатель: %d\n", aValue)

	// Получаем указатель на поле B
	bOffset := unsafe.Offsetof(s.B)
	bPtr := unsafe.Add(structPtr, bOffset)
	bValue := *(*int64)(bPtr)
	fmt.Printf("Поле B через указатель: %d\n", bValue)

	// Получаем указатель на поле C
	cOffset := unsafe.Offsetof(s.C)
	cPtr := unsafe.Add(structPtr, cOffset)
	cValue := *(*int32)(cPtr)
	fmt.Printf("Поле C через указатель: %d\n", cValue)

	// Изменяем значение через unsafe.Pointer
	*(*int32)(aPtr) = 999
	fmt.Printf("Структура после изменения A: %+v\n", s)
}
