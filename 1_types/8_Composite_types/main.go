package main

import "fmt"

// Объявление структуры на уровне пакета
type Person struct {
	Name string
	Age  int
}

// Объявление интерфейса на уровне пакета
type Speaker interface {
	Speak() string
}

// Метод для структуры Person, чтобы реализовать интерфейс Speaker
func (p Person) Speak() string {
	return fmt.Sprintf("Привет, меня зовут %s и мне %d лет.", p.Name, p.Age)
}

func main() {
	var arr1 [3]int = [3]int{1, 2, 3} // Полное объявление с типом
	arr2 := [3]int{1, 2, 3}           // Краткое объявление
	arr3 := [...]int{1, 2, 3, 4, 5}   // Компилятор сам вычислит длину (5)
	var arr4 [5]int                   // Объявление без инициализации (нулевой массив: [0,0,0,0,0])

	fmt.Println("--- Массивы ---")
	fmt.Println("arr1:", arr1)
	fmt.Println("arr2:", arr2)
	fmt.Println("arr3:", arr3)
	fmt.Println("arr4 (пустой):", arr4)

	var slice1 []int = []int{1, 2, 3} // Полное объявление
	slice2 := []int{1, 2, 3}          // Краткое объявление
	slice3 := make([]int, 3, 5)       // make() - длина 3, вместимость 5. Заполнено нулями: [0, 0, 0]
	slice4 := make([]int, 3)          // make() без явного указания вместимости
	var slice5 []int                  // Объявление без инициализации (nil слайс)

	slice2 = append(slice2, 4) // Добавление элемента в слайс
	slice3[0] = 10             // Изменение элемента

	fmt.Println("\n--- Слайсы ---")
	fmt.Println("slice1:", slice1, "len:", len(slice1), "cap:", cap(slice1))
	fmt.Println("slice2 (с добавленным элементом):", slice2)
	fmt.Println("slice3 (с измененным элементом):", slice3)
	fmt.Println("slice4:", slice4)
	fmt.Println("slice5 (nil):", slice5, "is nil:", slice5 == nil)

	var map1 map[string]int = map[string]int{"a": 1, "b": 2} // Полное объявление
	map2 := map[string]int{"a": 1, "b": 2}                   // Краткое объявление
	map3 := make(map[string]int)                             // make() - пустая, но инициализированная карта
	var map4 map[string]int                                  // Объявление без инициализации (nil карта)

	map3["c"] = 3 // Добавление элемента в инициализированную карту
	// map4["d"] = 4                                         // ОШИБКА! Запись в nil карту вызовет panic.

	// Проверка наличия ключа (запятая-ок паттерн)
	val, ok := map2["a"]
	if !ok {

	}
	fmt.Println("\n--- Карты ---")
	fmt.Println("map1:", map1)
	fmt.Println("map2:", map2)
	fmt.Println("map3:", map3)
	fmt.Println("map4 (nil):", map4, "is nil:", map4 == nil)
	fmt.Printf("Ключ 'a' в map2: значение=%d, найден=%t\n", val, ok)

	var p1 Person = Person{Name: "Alice", Age: 30} // Полное объявление с типом
	p2 := Person{Name: "Bob", Age: 25}             // Краткое объявление
	p3 := new(Person)                              // new() - вернет указатель на структуру с нулевыми полями
	var p4 Person                                  // Объявление без инициализации (пустая структура)
	p5 := Person{Age: 40}                          // Инициализация только части полей (Name будет "")
	p6 := struct {
		Name string
		Age  int
	}{"Inline", 99} // Анонимная структура (объявлена прямо в коде)

	p3.Name = "Charlie" // Изменение поля через указатель
	p4.Name = "Dave"

	fmt.Println("\n--- Структуры ---")
	fmt.Println("p1:", p1)
	fmt.Println("p2:", p2)
	fmt.Println("p3 (через new):", *p3)
	fmt.Println("p4 (пустая, затем измененная):", p4)
	fmt.Println("p5 (частичная инициализация):", p5)
	fmt.Println("p6 (анонимная):", p6)

	var ptr1 *int = new(int) // new() создает переменную и возвращает указатель на нее
	ptr2 := &p2              // Указатель на структуру p2
	var ptr3 *int            // Объявление без инициализации (nil указатель)
	x := 42
	ptr4 := &x // Указатель на переменную x
	// TODO уточнить
	var ptr5 = new(42)

	*ptr1 = 100 // Изменение значения по адресу

	fmt.Println("\n--- Указатели ---")
	fmt.Println("ptr1 (указатель на int, значение):", *ptr1)
	fmt.Println("ptr2 (указатель на структуру p2):", *ptr2)
	fmt.Println("ptr3 (nil):", ptr3, "is nil:", ptr3 == nil)
	fmt.Println("ptr4 (указатель на x):", *ptr4)
	fmt.Println("ptr5 (указатель на x):", *ptr5)

	chan1 := make(chan int)    // Небуферизированный канал
	chan2 := make(chan int, 5) // Буферизированный канал (емкость 5)
	var chan3 chan int         // Объявление без инициализации (nil канал)
	var chan4 <-chan int       // Канал "только для чтения" (read-only)
	var chan5 chan<- string    // Канал "только для записи" (write-only)

	// Асинхронная запись в буферизированный канал, чтобы избежать deadlock
	go func() {
		chan1 <- 42
		chan2 <- 1
		chan2 <- 2
	}()

	valChan := <-chan1 // Чтение из канала
	v2 := <-chan2

	fmt.Println("\n--- Каналы ---")
	fmt.Println("chan1 (небуферизированный, прочитано):", valChan)
	fmt.Println("chan2 (буферизированный, прочитано):", v2)
	fmt.Println("chan3 (nil):", chan3 == nil)
	fmt.Println("chan4 (read-only тип):", chan4)
	fmt.Println("chan5 (write-only тип):", chan5)

	var iface1 Speaker = p1          // p1 реализует Speaker, присваиваем интерфейсу
	var iface2 interface{} = "Hello" // Пустой интерфейс (может хранить ЛЮБОЙ тип)
	var iface3 interface{}           // Объявление без инициализации (nil интерфейс)

	// Преобразование типа (Type Assertion)
	if str, ok := iface2.(string); ok {
		fmt.Println("\n--- Интерфейсы ---")
		fmt.Println("iface1 (вызов метода):", iface1.Speak())
		fmt.Printf("iface2 (пустой интерфейс со строкой): значение='%s', тип=%T\n", str, iface2)
		fmt.Println("iface3 (nil):", iface3 == nil)
	}

	var fn1 func(int, int) int = func(a, b int) int { // Полное объявление
		return a + b
	}

	fn2 := func(s string) { // Краткое объявление
		fmt.Println("Функция вывела:", s)
	}
	var fn3 func() // Объявление без инициализации (nil функция)

	fmt.Println("\n--- Функции ---")
	fmt.Println("fn1 (сумма 5+10):", fn1(5, 10))
	fn2("Привет")
	fmt.Println("fn3 (nil):", fn3 == nil)
	// fn3() // ОШИБКА! Вызов nil функции вызовет panic.
}
