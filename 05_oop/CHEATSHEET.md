# Шпаргалка: ООП в Golang

**Быстрый справочник по всем ключевым концепциям**

---

## 1. Структуры (Struct)

### Определение
```go
type Point struct {
    X, Y float64
}

type User struct {
    ID   int64
    Name string
    Age  int
}
```

### Инициализация
```go
// Литерал с именами (рекомендуется)
p := Point{X: 1.0, Y: 2.0}
p := Point{X: 1.0}  // Y = 0 (zero value)

// Конструктор (с валидацией)
func NewUser(name string, age int) (*User, error) {
    if age < 0 {
        return nil, errors.New("invalid age")
    }
    return &User{Name: name, Age: age}, nil
}
```

### Zero value
```go
var p Point  // Point{X: 0, Y: 0}
var s string // ""
var i int    // 0
```

---

## 2. Методы (Methods)

### Определение
```go
// Value receiver (копия)
func (p Point) Distance() float64 {
    return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

// Pointer receiver (оригинал)
func (p *Point) Move(dx, dy float64) {
    p.X += dx
    p.Y += dy
}
```

### Выбор receiver'а

| Критерий | Value | Pointer |
|----------|-------|---------|
| Размер > 100 байт | ❌ | ✅ |
| Мутирует | ❌ | ✅ |
| Правило | Value | Pointer |
| Все методы | Один стиль! | |

```go
// ✅ Правильно: все pointer receivers
func (a *Account) Balance() float64 { ... }
func (a *Account) Deposit(amount float64) error { ... }

// ❌ Неправильно: смесь
func (a *Account) Balance() float64 { ... }      // value
func (a *Account) Deposit(amount float64) { ... }  // pointer
```

---

## 3. Инкапсуляция (Encapsulation)

### Видимость
```go
type User struct {
    Name   string  // Экспортируемо (заглавная)
    email  string  // Неэкспортируемо (строчная)
}
```

### Защита данных
```go
type Account struct {
    balance float64  // Приватное поле
}

func (a *Account) Balance() float64 {
    return a.balance  // Getter
}

func (a *Account) Deposit(amount float64) error {
    if amount <= 0 {
        return errors.New("must be positive")
    }
    a.balance += amount
    return nil  // Setter с валидацией
}
```

### Соглашение имен
```go
func (u *User) Name() string { ... }           // Getter
func (u *User) SetName(n string) error { ... } // Setter
func (a *Account) IsActive() bool { ... }      // Boolean
func (u *User) HasPermission(p string) bool { ... }  // Has-
```

---

## 4. Композиция (Composition)

### Без embedding (явная)
```go
type Car struct {
    engine Engine  // Явное поле
    color  string
}

car := Car{engine: Engine{}, color: "red"}
car.engine.Start()  // Явное обращение
```

### С embedding (встраивание)
```go
type Car struct {
    Engine  // Встроенное поле (anonymous)
    color   string
}

car := Car{Engine: Engine{}, color: "red"}
car.Start()  // Методы "поднимаются" (promoted)
car.Power    // Поля "поднимаются" (promoted)
```

### Разрешение конфликтов
```go
type Base struct{}
func (Base) Print() { fmt.Println("Base") }

type Derived struct{ Base }
func (Derived) Print() { fmt.Println("Derived") }

d := Derived{}
d.Print()        // "Derived" (метод Derived скрывает Base)
d.Base.Print()   // "Base" (явно обращаемся к Base)
```

---

## 5. Интерфейсы (Interfaces)

### Определение
```go
type Reader interface {
    Read([]byte) (int, error)
}

type Writer interface {
    Write([]byte) (int, error)
}

type ReadWriter interface {
    Reader
    Writer
}
```

### Реализация (неявная)
```go
// MyReader реализует Reader, если имеет метод Read()
type MyReader struct{}

func (mr MyReader) Read(b []byte) (int, error) {
    // Реализация
}

var r Reader = MyReader{}  // ✅ Автоматически
```

### Acceptance
```go
// Принимайте интерфейсы
func Copy(dst Writer, src Reader) error {
    // Работает с любым Reader/Writer
}

// Возвращайте конкретные типы
func Open(name string) (*os.File, error) {
    return os.Open(name)
}
```

---

## 6. Выравнивание памяти (Memory Alignment)

### Проблема
```go
// ❌ 16 bytes (много padding)
type Bad struct {
    flag bool      // 1 byte + 7 padding
    value int64    // 8 bytes
}

// ✅ 16 bytes (эффективнее)
type Good struct {
    value int64    // 8 bytes
    flag bool      // 1 byte
}
```

### Оптимизация
```go
// Правило: от большего к меньшему
type Optimized struct {
    field1 int64    // 8 bytes
    field2 int32    // 4 bytes
    field3 int16    // 2 bytes
    field4 bool     // 1 byte
}
```

### Проверка
```go
unsafe.Sizeof(Bad{})              // 16
unsafe.Offsetof(Bad{}.flag)       // 0
unsafe.Offsetof(Bad{}.value)      // 8
```

---

## 7. Обработка ошибок (Error Handling)

### Sentinel errors
```go
var (
    ErrNotFound = errors.New("not found")
    ErrTimeout  = errors.New("timeout")
)

if err == ErrNotFound { ... }
```

### Error types
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("%s: %s", e.Field, e.Message)
}

if err := someFunc(); err != nil {
    var ve *ValidationError
    if errors.As(err, &ve) {
        fmt.Printf("Validation failed on %s\n", ve.Field)
    }
}
```

### Оборачивание ошибок
```go
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Теперь можно распаковать
if errors.Is(err, os.ErrNotExist) { ... }
```

---

## 8. Паттерны проектирования (Design Patterns)

### Factory (NewXxx)
```go
func NewUser(name string) (*User, error) {
    if name == "" {
        return nil, errors.New("name required")
    }
    return &User{name: name}, nil
}
```

### Repository
```go
type UserRepository interface {
    GetByID(id int64) (*User, error)
    Save(u *User) error
}

type Service struct {
    repo UserRepository
}
```

### Strategy
```go
type PaymentStrategy interface {
    Pay(amount float64) error
}

type Order struct {
    items    []*Item
    strategy PaymentStrategy
}

func (o *Order) Checkout() error {
    return o.strategy.Pay(o.Total())
}
```

### Builder
```go
query := NewQueryBuilder().
    Select("id", "name").
    From("users").
    Where("age > 18").
    Build()
```

### Decorator
```go
type LoggingStore struct {
    Store  // Embedded interface
}

func (ls *LoggingStore) Get(key string) (string, error) {
    log.Printf("Getting %s", key)
    return ls.Store.Get(key)
}
```

---

## 9. Тестирование (Testing)

### Unit test
```go
func TestDeposit(t *testing.T) {
    acc := NewAccount(1000)
    err := acc.Deposit(500)
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    
    if acc.Balance() != 1500 {
        t.Errorf("expected 1500, got %v", acc.Balance())
    }
}
```

### Table-driven test
```go
func TestDiscount(t *testing.T) {
    tests := []struct {
        name     string
        amount   float64
        percent  int
        expected float64
        wantErr  bool
    }{
        {"valid", 100, 10, 90, false},
        {"negative", -100, 10, 0, true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Discount(tt.amount, tt.percent)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
            }
            if got != tt.expected {
                t.Errorf("got %v, expected %v", got, tt.expected)
            }
        })
    }
}
```

### Mocking
```go
type MockRepository struct {
    GetFunc func(id int64) (*User, error)
}

func (m *MockRepository) Get(id int64) (*User, error) {
    return m.GetFunc(id)
}
```

---

## 10. Быстрая справка по типам

### Встроенные типы
```go
bool                  // true, false
int, int8, int16, ... // целые числа
float32, float64      // вещественные
string                // текст
complex64, complex128 // комплексные числа
error                 // интерфейс ошибки
interface{}           // любой тип
```

### Размеры в памяти
```go
unsafe.Sizeof(bool{})      // 1
unsafe.Sizeof(int8{})      // 1
unsafe.Sizeof(int16{})     // 2
unsafe.Sizeof(int32{})     // 4
unsafe.Sizeof(int64{})     // 8
unsafe.Sizeof(float64{})   // 8
unsafe.Sizeof(string{})    // 16
unsafe.Sizeof([]byte{})    // 24
unsafe.Sizeof(map[string]int{})  // 8
```

---

## 11. Полезные команды

```bash
# Форматирование
gofmt -w ./...
goimports -w ./...

# Статический анализ
go vet ./...
golangci-lint run ./...

# Тестирование
go test -v ./...
go test -cover ./...

# Сборка
go build ./...

# Документация
godoc -http=:8080

# Профилирование памяти
go run -memprofile=mem.prof ./...
go tool pprof mem.prof

# Размеры структур
structlayout ./mypackage
```

---

## 12. Checklist для code review

- [ ] Правильный выбор value/pointer receiver
- [ ] Все методы одного типа используют один стиль
- [ ] Приватные поля используют setter'ы с валидацией
- [ ] Экспортируются интерфейсы, не конкретные типы
- [ ] Ошибки обработаны и обёрнуты
- [ ] Нет игнорирования ошибок (`_ =`)
- [ ] Структуры оптимизированы по памяти
- [ ] Есть конструкторы с валидацией
- [ ] Интерфейсы маленькие и сфокусированные
- [ ] Есть тесты (минимум 80% покрытие)

---

## 13. Типичные ошибки

```go
// ❌ ОШИБКА 1: Смешивание receivers
func (u User) GetName() string { ... }
func (u *User) SetName(name string) { ... }

// ❌ ОШИБКА 2: Игнорирование ошибок
_ = account.Withdraw(100)

// ❌ ОШИБКА 3: Экспортирование типов вместо интерфейсов
type userStore struct { ... }  // ❌
type UserStore interface { ... }  // ✅

// ❌ ОШИБКА 4: Неправильное выравнивание
struct { flag bool; value int64 }  // 16 bytes
struct { value int64; flag bool }  // 16 bytes, но эффективнее

// ❌ ОШИБКА 5: Embedding вместо композиции
struct { Logger; config string }  // Logger IS-A? Нет!
struct { logger *Logger; config string }  // HAS-A ✅
```

---

## 14. Ресурсы

### Документация
- https://go.dev/doc/ — официальная документация
- https://go.dev/ref/spec — спецификация
- https://gobyexample.com/ — примеры

### Книги
- "The Go Programming Language" (Donovan, Kernighan)
- "Concurrency in Go" (Cox-Buday)

### Стайл гайды
- https://github.com/golang/go/wiki/CodeReviewComments
- https://github.com/uber-go/guide

---

**Сохраните эту шпаргалку как закладку!** 📌

Вернитесь сюда, когда нужен быстрый справочник.
