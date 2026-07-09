// Package main - Практические примеры из лекции ООП в Golang
package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strings"
	"sync"
	"unsafe"
)

// ============================================================================
// Раздел 2: СТРУКТУРЫ
// ============================================================================

// Point - простая структура с двумя полями
type Point struct {
	X, Y float64
}

// Example 2.1: Инициализация структур
func Example_InitializationStructs() {
	fmt.Println("=== 2.1: Инициализация структур ===")

	// Способ 1: Литерал без имен (опасный)
	p1 := Point{3.0, 4.0}
	fmt.Printf("p1 (positional): %+v\n", p1)

	// Способ 2: Литерал с именами полей (рекомендуется)
	p2 := Point{X: 5.0, Y: 12.0}
	fmt.Printf("p2 (named fields): %+v\n", p2)

	// Способ 3: Частичная инициализация
	p3 := Point{X: 7.0} // Y будет 0
	fmt.Printf("p3 (partial): %+v\n", p3)

	// Адресс структуры
	ptr := &p2
	fmt.Printf("ptr.X = %v\n", ptr.X)
	ptr.X = 100 // Автоматическое разыменование
	fmt.Printf("После изменения через указатель: %+v\n", p2)

	fmt.Println()
}

// Example 2.2: Zero value
func Example_ZeroValue() {
	fmt.Println("=== 2.2: Zero value (нулевые значения) ===")

	type Person struct {
		Name   string
		Age    int
		Salary float64
		Active bool
	}

	var p Person // Не инициализирована
	fmt.Printf("Zero value Person: %+v\n", p)
	fmt.Printf("  Name (string): %q\n", p.Name)
	fmt.Printf("  Age (int): %d\n", p.Age)
	fmt.Printf("  Salary (float64): %v\n", p.Salary)
	fmt.Printf("  Active (bool): %v\n", p.Active)

	fmt.Println()
}

// Example 2.3: Сравнение структур
func Example_CompareStructs() {
	fmt.Println("=== 2.3: Сравнение структур ===")

	type Color struct {
		R, G, B uint8
	}

	c1 := Color{255, 0, 0}
	c2 := Color{255, 0, 0}
	c3 := Color{0, 255, 0}

	fmt.Printf("c1 = %v\n", c1)
	fmt.Printf("c2 = %v\n", c2)
	fmt.Printf("c3 = %v\n", c3)
	fmt.Printf("c1 == c2: %v\n", c1 == c2)
	fmt.Printf("c1 == c3: %v\n", c1 == c3)

	fmt.Println()
}

// ============================================================================
// Раздел 3: ИНКАПСУЛЯЦИЯ
// ============================================================================

// Account - счёт с защитой данных
type Account struct {
	balance float64 // Неэкспортируемое поле
	owner   string
}

// NewAccount - конструктор с валидацией
func NewAccount(owner string, initialBalance float64) *Account {
	if initialBalance < 0 {
		initialBalance = 0
	}
	return &Account{
		balance: initialBalance,
		owner:   owner,
	}
}

// Balance - getter
func (a *Account) Balance() float64 {
	return a.balance
}

// Deposit - пополнение со счёта
func (a *Account) Deposit(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("deposit amount must be positive: %v", amount)
	}
	a.balance += amount
	return nil
}

// Withdraw - снятие со счёта
func (a *Account) Withdraw(amount float64) error {
	if amount <= 0 {
		return fmt.Errorf("withdraw amount must be positive: %v", amount)
	}
	if a.balance < amount {
		return fmt.Errorf("insufficient funds: have %v, need %v", a.balance, amount)
	}
	a.balance -= amount
	return nil
}

// Example 3.1: Инкапсуляция и защита данных
func Example_Encapsulation() {
	fmt.Println("=== 3: Инкапсуляция и защита данных ===")

	acc := NewAccount("Alice", 1000)
	fmt.Printf("Initial balance: %v\n", acc.Balance())

	// Правильное использование
	if err := acc.Deposit(500); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("After deposit 500: %v\n", acc.Balance())

	// Попытка некорректной операции
	if err := acc.Withdraw(2000); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	// Попытка некорректной операции
	if err := acc.Deposit(-100); err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	fmt.Println()
}

// ============================================================================
// Раздел 4: МЕТОДЫ И RECEIVERS
// ============================================================================

// Counter - простой счётчик для демонстрации receivers
type Counter struct {
	count int
}

// Value receiver - получает копию структуры
func (c Counter) Increment_BadExample() {
	c.count++ // Изменяет КОПИЮ!
}

// Pointer receiver - получает указатель
func (c *Counter) Increment() {
	c.count++ // Изменяет оригинал
}

// Example 4.1: Value vs Pointer receivers
func Example_ValueVsPointerReceivers() {
	fmt.Println("=== 4.1: Value vs Pointer receivers ===")

	counter := Counter{count: 0}

	// Value receiver - НЕ изменяет оригинал
	counter.Increment_BadExample()
	fmt.Printf("After Increment_BadExample(): %d (не изменилось!)\n", counter.count)

	// Pointer receiver - изменяет оригинал
	counter.Increment()
	fmt.Printf("After Increment(): %d\n", counter.count)

	counter.Increment()
	fmt.Printf("After Increment(): %d\n", counter.count)

	fmt.Println()
}

// Kilometers - тип на основе встроенного типа
type Kilometers float64

func (k Kilometers) Miles() float64 {
	return float64(k) * 0.621371
}

func (k Kilometers) String() string {
	return fmt.Sprintf("%.2f km", k)
}

// Example 4.2: Методы на пользовательских типах
func Example_MethodsOnCustomTypes() {
	fmt.Println("=== 4.2: Методы на пользовательских типах ===")

	distance := Kilometers(5)
	fmt.Printf("Distance: %s\n", distance)
	fmt.Printf("In miles: %.2f\n", distance.Miles())

	fmt.Println()
}

// ============================================================================
// Раздел 5: КОМПОЗИЦИЯ И EMBEDDING
// ============================================================================

// Engine - двигатель
type Engine struct {
	Power int // Лошадиные силы
	Fuel  int // Литры топлива
}

func (e *Engine) Start() string {
	if e.Fuel <= 0 {
		return "Cannot start: no fuel"
	}
	return fmt.Sprintf("Engine started with %d HP", e.Power)
}

// Vehicle - транспортное средство (явная композиция)
type Vehicle struct {
	Engine Engine
	Color  string
}

// Car - автомобиль (через embedding)
type Car struct {
	Engine // Embedded field - методы и поля поднимаются (promoted)
	Color  string
	Wheels int
}

// Example 5.1: Композиция vs Embedding
func Example_CompositionEmbedding() {
	fmt.Println("=== 5: Композиция и Embedding ===")

	// Явная композиция
	v := Vehicle{
		Engine: Engine{Power: 150, Fuel: 50},
		Color:  "blue",
	}
	fmt.Printf("Vehicle: %+v\n", v)
	fmt.Printf("Vehicle engine power: %d\n", v.Engine.Power)
	fmt.Printf("Vehicle start: %s\n", v.Engine.Start())

	// Embedding
	car := Car{
		Engine: Engine{Power: 200, Fuel: 60},
		Color:  "red",
		Wheels: 4,
	}
	fmt.Printf("\nCar: %+v\n", car)
	fmt.Printf("Car power (promoted): %d\n", car.Power)   // ✅ Promoted field
	fmt.Printf("Car start (promoted): %s\n", car.Start()) // ✅ Promoted method
	fmt.Printf("Car wheels: %d\n", car.Wheels)

	fmt.Println()
}

// Reader интерфейс
type Reader interface {
	Read() string
}

// StringReader реализует Reader
type StringReader struct {
	data string
}

func (sr *StringReader) Read() string {
	return sr.data
}

// Logger обогащает Reader логированием
type LoggedReader struct {
	Reader // Embedded interface
	log    *log.Logger
}

func (lr *LoggedReader) Read() string {
	lr.log.Println("Reading data...")
	return lr.Reader.Read()
}

// Example 5.2: Embedding интерфейсов для декораторов
func Example_EmbeddingInterfaces() {
	fmt.Println("=== 5.2: Embedding интерфейсов (Decorator pattern) ===")

	reader := &StringReader{data: "Hello, World!"}

	logger := log.New(log.Writer(), "[LOG] ", log.LstdFlags)
	logged := &LoggedReader{Reader: reader, log: logger}

	result := logged.Read()
	fmt.Printf("Result: %s\n", result)

	fmt.Println()
}

// ============================================================================
// Раздел 6: ВЫРАВНИВАНИЕ ПАМЯТИ
// ============================================================================

// BadAlignment - плохое выравнивание (много padding)
type BadAlignment struct {
	flag  bool  // 1 байт + 7 padding
	value int64 // 8 байт (требует выравнивания на 8)
}

// GoodAlignment - хорошее выравнивание
type GoodAlignment struct {
	value int64 // 8 байт
	flag  bool  // 1 байт
}

// Example 6: Выравнивание памяти
func Example_MemoryAlignment() {
	fmt.Println("=== 6: Выравнивание памяти ===")

	bad := BadAlignment{}
	good := GoodAlignment{}

	fmt.Printf("BadAlignment size: %d bytes\n", unsafe.Sizeof(bad))
	fmt.Printf("BadAlignment.flag offset: %d\n", unsafe.Offsetof(bad.flag))
	fmt.Printf("BadAlignment.value offset: %d\n", unsafe.Offsetof(bad.value))

	fmt.Printf("\nGoodAlignment size: %d bytes\n", unsafe.Sizeof(good))
	fmt.Printf("GoodAlignment.value offset: %d\n", unsafe.Offsetof(good.value))
	fmt.Printf("GoodAlignment.flag offset: %d\n", unsafe.Offsetof(good.flag))

	fmt.Println()
}

// ============================================================================
// Раздел 7: BEST PRACTICES
// ============================================================================

// ValidationError - кастомная ошибка
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation error on field %q: %s", e.Field, e.Message)
}

// User - пользователь с валидацией
type User struct {
	id    int64
	name  string
	email string
	age   int
}

// NewUser создает нового пользователя
func NewUser(name, email string) (*User, error) {
	if name == "" || email == "" {
		return nil, errors.New("name and email are required")
	}
	return &User{
		name:  name,
		email: email,
		age:   0,
	}, nil
}

// SetAge устанавливает возраст с валидацией
func (u *User) SetAge(age int) error {
	if age < 0 || age > 150 {
		return &ValidationError{
			Field:   "age",
			Message: "must be between 0 and 150",
		}
	}
	u.age = age
	return nil
}

// Example 7.1: Best Practices - обработка ошибок
func Example_ErrorHandling() {
	fmt.Println("=== 7.1: Best Practices - обработка ошибок ===")

	user, err := NewUser("Alice", "alice@example.com")
	if err != nil {
		log.Fatal(err)
	}

	// Успешное установление возраста
	if err := user.SetAge(30); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("User age set successfully: %d\n", user.age)

	// Ошибочная валидация
	if err := user.SetAge(200); err != nil {
		var ve *ValidationError
		if errors.As(err, &ve) {
			fmt.Printf("Validation failed on %q: %s\n", ve.Field, ve.Message)
		}
	}

	fmt.Println()
}

// Store интерфейс
type Store interface {
	Get(key string) (string, error)
	Set(key, value string) error
}

// MemoryStore реализует Store
type MemoryStore struct {
	data map[string]string
	mu   sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		data: make(map[string]string),
	}
}

func (ms *MemoryStore) Get(key string) (string, error) {
	ms.mu.RLock()
	defer ms.mu.RUnlock()
	value, ok := ms.data[key]
	if !ok {
		return "", errors.New("key not found")
	}
	return value, nil
}

func (ms *MemoryStore) Set(key, value string) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()
	ms.data[key] = value
	return nil
}

// LoggingStore - декоратор с логированием
type LoggingStore struct {
	Store // Embedded interface
	log   *log.Logger
}

func (ls *LoggingStore) Get(key string) (string, error) {
	ls.log.Printf("Getting key: %q", key)
	value, err := ls.Store.Get(key)
	if err != nil {
		ls.log.Printf("Error: %v", err)
	}
	return value, err
}

func (ls *LoggingStore) Set(key, value string) error {
	ls.log.Printf("Setting key: %q to %q", key, value)
	return ls.Store.Set(key, value)
}

// Example 7.2: Паттерны проектирования - Decorator
func Example_DecoratorPattern() {
	fmt.Println("=== 7.2: Паттерны проектирования - Decorator ===")

	baseStore := NewMemoryStore()
	logger := log.New(log.Writer(), "[STORE] ", log.LstdFlags)
	loggedStore := &LoggingStore{Store: baseStore, log: logger}

	// Все операции будут логироваться
	loggedStore.Set("name", "Alice")
	loggedStore.Set("age", "30")

	value, err := loggedStore.Get("name")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Got value: %q\n", value)
	}

	fmt.Println()
}

// QueryBuilder - паттерн Builder
type QueryBuilder struct {
	selects []string
	from    string
	wheres  []string
}

func NewQueryBuilder() *QueryBuilder {
	return &QueryBuilder{
		selects: make([]string, 0), // Пустой список полей
		from:    "",                // Пустая строка для таблицы
		wheres:  make([]string, 0), // Пустой список условий
	}
}

func (qb *QueryBuilder) Select(fields ...string) *QueryBuilder {
	qb.selects = append(qb.selects, fields...)
	return qb
}

func (qb *QueryBuilder) From(table string) *QueryBuilder {
	qb.from = table
	return qb
}

func (qb *QueryBuilder) Where(condition string) *QueryBuilder {
	qb.wheres = append(qb.wheres, condition)
	return qb
}

func (qb *QueryBuilder) Build() string {
	var sb strings.Builder

	sb.WriteString("SELECT ")
	sb.WriteString(strings.Join(qb.selects, ", "))

	if qb.from != "" {
		sb.WriteString(" FROM ")
		sb.WriteString(qb.from)
	}

	if len(qb.wheres) > 0 {
		sb.WriteString(" WHERE ")
		sb.WriteString(strings.Join(qb.wheres, " AND "))
	}

	return sb.String()
}

// Example 7.3: Паттерны проектирования - Builder
func Example_BuilderPattern() {
	fmt.Println("=== 7.3: Паттерны проектирования - Builder ===")

	query := NewQueryBuilder().
		Select("id", "name", "email").
		From("users").
		Where("age > 18").
		Build()

	fmt.Printf("Query: %s\n", query)

	fmt.Println()
}

// ============================================================================
// Shape интерфейс - практическое задание
// ============================================================================

type Shape interface {
	Area() float64
	Perimeter() float64
}

type Rectangle struct {
	width, height float64
}

func (r Rectangle) Area() float64 {
	return r.width * r.height
}

func (r Rectangle) Perimeter() float64 {
	return 2 * (r.width + r.height)
}

type Circle struct {
	radius float64
}

func (c Circle) Area() float64 {
	return math.Pi * c.radius * c.radius
}

func (c Circle) Perimeter() float64 {
	return 2 * math.Pi * c.radius
}

// Example 8: Практическое задание - Фигуры
func Example_PracticalTask() {
	fmt.Println("=== 8: Практическое задание - Геометрические фигуры ===")

	shapes := []Shape{
		Rectangle{width: 3, height: 4},
		Circle{radius: 5},
		Rectangle{width: 2, height: 8},
	}

	for i, shape := range shapes {
		fmt.Printf("Shape %d:\n", i+1)
		fmt.Printf("  Area: %.2f\n", shape.Area())
		fmt.Printf("  Perimeter: %.2f\n", shape.Perimeter())
	}

	fmt.Println()
}

// ============================================================================
// MAIN - запуск всех примеров
// ============================================================================

func main() {
	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║  ООП в Golang - Практические примеры из лекции            ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
	fmt.Println()

	// Раздел 2: Структуры
	Example_InitializationStructs()
	Example_ZeroValue()
	Example_CompareStructs()

	// Раздел 3: Инкапсуляция
	Example_Encapsulation()

	// Раздел 4: Методы и Receivers
	Example_ValueVsPointerReceivers()
	Example_MethodsOnCustomTypes()

	// Раздел 5: Композиция и Embedding
	Example_CompositionEmbedding()
	Example_EmbeddingInterfaces()

	// Раздел 6: Выравнивание памяти
	Example_MemoryAlignment()

	// Раздел 7: Best Practices
	Example_ErrorHandling()
	Example_DecoratorPattern()
	Example_BuilderPattern()
	//
	// Раздел 8: Практическое задание
	Example_PracticalTask()

	fmt.Println("╔════════════════════════════════════════════════════════════╗")
	fmt.Println("║  Все примеры выполнены успешно!                           ║")
	fmt.Println("╚════════════════════════════════════════════════════════════╝")
}
