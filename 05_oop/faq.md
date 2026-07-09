# FAQ: ООП в Golang

**Часто задаваемые вопросы и ответы по ООП в Go**

---

## Структуры и типы

### Q: В Go нет классов. Это означает, что нет ООП?

**A:** Нет. Go имеет ООП, но это функциональный ООП, а не классовый.

Go предоставляет инструменты для ООП:
- **Struct** для группировки данных
- **Methods** для поведения
- **Interfaces** для полиморфизма
- **Embedding** для композиции

```go
// ООП в Go работает через методы и интерфейсы
type Animal interface {
    Speak() string
}

type Dog struct {
    name string
}

func (d Dog) Speak() string {
    return "Woof"
}

// Полиморфизм
var animal Animal = Dog{name: "Rex"}
fmt.Println(animal.Speak())  // "Woof"
```

---

### Q: Как определить "приватные" поля?

**A:** Используйте строчные буквы. Поля с заглавной буквы экспортируются, со строчной — нет.

```go
type User struct {
    Name     string  // Экспортируемое (видно из других пакетов)
    password string  // Неэкспортируемое (видно только в пакете)
}

// В другом пакете:
user := main.User{Name: "Alice"}         // OK
user := main.User{password: "secret"}    // Ошибка компиляции
```

---

### Q: Зачем нужны конструкторы, если можно использовать литерал?

**A:** Конструкторы позволяют:
1. **Валидировать** данные
2. **Инициализировать** сложные поля
3. **Добавлять логику** при создании
4. **Изменять реализацию** без изменения API

```go
// Опасный литерал
user := User{Name: "", Age: -5}  // Невалидные данные!

// Безопасный конструктор
user, err := NewUser("Alice", 30)  // Валидация внутри
```

---

## Методы и Receivers

### Q: Когда использовать value receiver, а когда pointer?

**A:** 

| Случай | Receiver | Причина |
|--------|----------|---------|
| **Маленькие структуры** | Value | Дешевое копирование (< 100 байт) |
| **Большие структуры** | Pointer | Избегаем дорогого копирования |
| **Мутирующие методы** | Pointer | Нужно изменить оригинал |
| **Неизменяемые методы** | Value | Логика read-only |

**Практическое правило:** если хоть один метод использует pointer receiver, все остальные методы типа должны тоже.

```go
// Последовательно
type Account struct { balance float64 }
func (a *Account) Deposit(amount float64) { ... }
func (a *Account) Withdraw(amount float64) { ... }
func (a *Account) Balance() float64 { ... }

// Непоследовательно (не рекомендуется)
func (a *Account) Deposit(amount float64) { ... }    // pointer
func (a Account) Balance() float64 { ... }           // value
```

---

### Q: Почему я могу вызывать методы с pointer receiver на value?

**A:** Go автоматически берёт адрес:

```go
counter := Counter{count: 0}
counter.Increment()  // Эквивалентно: (&counter).Increment()

// Это работает, только если:
// 1. counter это переменная (не результат выражения)
// 2. метод имеет pointer receiver
```

**Но:**

```go
var counts = []Counter{{0}, {1}, {2}}
counts[0].Increment()  // Ошибка! Нельзя взять адрес элемента слайса в данном контексте
```

---

### Q: Могу ли я определить методы на встроенных типах?

**A:** Нет напрямую, но можно создать новый тип:

```go
// Не работает
func (i int) Double() int { return i * 2 }

// Работает
type MyInt int

func (i MyInt) Double() MyInt { return i * 2 }

n := MyInt(5)
fmt.Println(n.Double())  // 10
```

---

## Инкапсуляция

### Q: Как защитить данные от неправильного использования?

**A:** Используйте неэкспортируемые поля + методы с валидацией:

```go
type BankAccount struct {
    balance float64  // Приватное поле
}

func (a *BankAccount) Deposit(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("amount must be positive")
    }
    a.balance += amount
    return nil
}

// Клиент не может сделать
account.balance = -999  // Ошибка компиляции
account.Deposit(-999)   // Ошибка во время выполнения
```

---

### Q: Нужно ли всегда использовать getter'ы и setter'ы?

**A:** Нет. Простые публичные поля часто OK:

```go
// Простые данные — публичные поля
type Point struct {
    X, Y float64
}

// Данные с инвариантами — getter'ы/setter'ы
type Score struct {
    value int
}

func (s *Score) Set(v int) error {
    if v < 0 || v > 100 {
        return errors.New("score must be 0-100")
    }
    s.value = v
    return nil
}
```

---

## Композиция и Embedding

### Q: В чём разница между композицией и embedding?

**A:** 

```go
// Композиция (явная)
type Car struct {
    engine Engine  // Явное поле
}

car := Car{engine: Engine{Power: 100}}
car.engine.Start()  // Нужно явно обращаться

// Embedding (встраивание)
type Car struct {
    Engine  // Встроенное поле (anonymous)
}

car := Car{Engine: Engine{Power: 100}}
car.Start()  // Методы "поднимаются" (promoted)
```

**Когда использовать:**
- **Композиция:** отношение "HAS-A" (Car HAS-A Engine)
- **Embedding:** отношение "IS-A" (Admin IS-A User)

---

### Q: Что происходит, если встроенный и встраивающий типы имеют метод с одинаковым именем?

**A:** Метод встраивающего типа "скрывает" метод встроенного типа:

```go
type Base struct{}
func (Base) Print() { fmt.Println("Base") }

type Derived struct{ Base }
func (Derived) Print() { fmt.Println("Derived") }

d := Derived{}
d.Print()        // "Derived" (вызывается метод Derived)
d.Base.Print()   // "Base" (явно обращаемся к Base)
```

---

### Q: Могу ли я встраивать несколько типов одновременно?

**A:** Да:

```go
type Reader interface {
    Read() string
}

type Writer interface {
    Write(string) error
}

type ReadWriter struct {
    Reader
    Writer  // Множественное встраивание
}

// ReadWriter имеет методы Read() и Write()
```

**Но:** избегайте конфликтов имён:

```go
// Конфликт
type A struct{ Name string }
type B struct{ Name string }
type C struct {
    A
    B
}

c := C{}
fmt.Println(c.Name)  // Какое Name? Амбигуально!
fmt.Println(c.A.Name)  // OK — явно
```

---

## Интерфейсы

### Q: Нужно ли явно реализовать интерфейс?

**A:** Нет. Go использует структурную типизацию (duck typing):

```go
type Reader interface {
    Read([]byte) (int, error)
}

type MyReader struct{}

func (mr MyReader) Read(b []byte) (int, error) {
    // Реализация
}

var r Reader = MyReader{}  // Автоматически реализует Reader
```

Если тип имеет все методы интерфейса, он реализует интерфейс.

---

### Q: Могу ли я получить список всех типов, которые реализуют интерфейс?

**A:** Нет. Go не отслеживает это во время компиляции. Это особенность структурной типизации.

Можно использовать рефлексию во время выполнения (но это медленно и сложно).

---

### Q: Что такое пустой интерфейс `interface{}`?

**A:** Это интерфейс без методов. Любой тип его реализует:

```go
var i interface{}  // Может быть любым значением

i = 42
i = "hello"
i = []int{1, 2, 3}

// Type assertion
str, ok := i.(string)  // Проверяем, что это string
```

В Go 1.18+ используйте `any` вместо `interface{}`:

```go
var i any = 42
```

---

## Выравнивание и память

### Q: Почему размер структуры больше, чем сумма размеров полей?

**A:** Это padding (заполнение) для выравнивания:

```go
type Bad struct {
    flag bool      // 1 byte + 7 padding
    value int64    // 8 bytes (требует выравнивания на 8)
}

unsafe.Sizeof(Bad{})  // 16 bytes, не 9
```

CPU читает память эффективнее, если данные выравнены.

---

### Q: Как оптимизировать размер структуры?

**A:** Группируйте поля от большего к меньшему размеру:

```go
// Неоптимально
type Person struct {
    Name   string      // 16 bytes
    Age    int8        // 1 byte
    Height float64     // 8 bytes
}

// Оптимально
type Person struct {
    Name   string      // 16 bytes
    Height float64     // 8 bytes
    Age    int8        // 1 byte
}
```

**Важно:** это микро-оптимизация. Обычно не критично.

---

### Q: Как измерить размер структуры?

**A:** Используйте `unsafe.Sizeof`:

```go
import "unsafe"

type User struct {
    id    int64
    name  string
}

fmt.Println(unsafe.Sizeof(User{}))              // 24 bytes
fmt.Println(unsafe.Offsetof(User{}.id))         // 0
fmt.Println(unsafe.Offsetof(User{}.name))       // 8
```

Для анализа структуры:
```bash
go install golang.org/x/exp/cmd/structlayout@latest
structlayout ./mypackage
```

---

## Обработка ошибок

### Q: Что лучше использовать: sentinel errors или error types?

**A:** 

| Тип | Примеры | Когда |
|-----|---------|-------|
| **Sentinel** | `ErrNotFound`, `EOF` | Простые ошибки без контекста |
| **Type** | `ValidationError`, `TimeoutError` | Ошибки с дополнительной информацией |
| **Both** | Комбинация обоих | Иерархия ошибок |

```go
// Sentinel error
var ErrNotFound = errors.New("not found")

if err == ErrNotFound { ... }

// Error type
type NotFoundError struct {
    ID string
}

func (e *NotFoundError) Error() string {
    return fmt.Sprintf("item %q not found", e.ID)
}

if _, ok := err.(*NotFoundError); ok { ... }
```

---

### Q: Всегда ли нужно оборачивать ошибки с `%w`?

**A:** Да, во время разработки. `%w` позволяет:

```go
// Теряем информацию об оригинальной ошибке
if err != nil {
    return fmt.Errorf("operation failed: %v", err)
}

// Сохраняем оригинальную ошибку
if err != nil {
    return fmt.Errorf("operation failed: %w", err)
}

// Теперь можно распаковать:
if errors.Is(err, os.ErrNotExist) { ... }
```

---

### Q: Нужно ли проверять ошибку сразу же?

**A:** Да, обычно. Исключение — когда вы точно знаете, что операция не может ошибиться:

```go
// Хорошо: проверяем каждую ошибку
data, err := os.ReadFile("file.txt")
if err != nil {
    return fmt.Errorf("failed to read file: %w", err)
}

// OK: знаем, что JSON.Unmarshal не может ошибиться
_ = json.Unmarshal(validJSON, &data)  // validJSON гарантированно валидна

// Плохо: игнорируем возможную ошибку
json.Unmarshal(userInput, &data)  // Может ошибиться!
```

---

## Тестирование

### Q: Как тестировать методы с pointer receiver?

**A:** Никак не отличается:

```go
type Counter struct{ count int }
func (c *Counter) Increment() { c.count++ }

func TestIncrement(t *testing.T) {
    c := &Counter{}
    c.Increment()
    if c.count != 1 {
        t.Errorf("expected 1, got %d", c.count)
    }
}
```

---

### Q: Как мокировать интерфейсы?

**A:** Создайте mock-тип, который реализует интерфейс:

```go
type MockDatabase struct {
    GetFunc func(id string) (string, error)
}

func (m *MockDatabase) Get(id string) (string, error) {
    return m.GetFunc(id)
}

func TestService(t *testing.T) {
    mock := &MockDatabase{
        GetFunc: func(id string) (string, error) {
            return "mocked data", nil
        },
    }
    
    svc := &Service{db: mock}
    result := svc.Process(context.Background(), "123")
    // Проверяем результат
}
```

Альтернатива: используйте стороннюю библиотеку вроде `github.com/golang/mock`.

---

## Performance

### Q: Как я знаю, использует ли мой метод value или pointer receiver?

**A:** Посмотрите сигнатуру:

```go
func (p Point) Distance() float64 { ... }        // Value receiver
func (p *Point) Move(dx, dy float64) { ... }     // Pointer receiver
```

Вызов идентичен для обеих, но семантика другая.

---

### Q: Создаёт ли Go копию структуры, если я передам её в функцию?

**A:** Да, если не использовать указатель:

```go
func modifyByValue(p Point) {       // Копия Point
    p.X = 999                       // Не влияет на оригинал
}

func modifyByPointer(p *Point) {    // Указатель
    p.X = 999                       // Влияет на оригинал
}

p := Point{X: 1, Y: 2}
modifyByValue(p)        // Копируется
modifyByPointer(&p)     // Адрес передается (дешево)
```

Для больших структур всегда передавайте указатель.

---

### Q: Нужны ли мне object pools для производительности?

**A:** Обычно нет. Сборка мусора Go достаточно эффективна.

Object pool может помочь, только если:
1. Вы профилировали и нашли проблему
2. Создаёте миллионы объектов в секунду
3. Это критичный код (hot path)

```go
import "sync"

pool := &sync.Pool{
    New: func() interface{} {
        return &ExpensiveObject{}
    },
}

obj := pool.Get().(*ExpensiveObject)
defer pool.Put(obj)
```

---

## Общие вопросы

### Q: Почему в Go нет классического наследования?

**A:** Потому что:

1. **Композиция проще и гибче** — вы явно выбираете, какие поведения хотите переиспользовать
2. **Избегаем Fragile Base Class problem** — изменение базового класса не сломает производные классы
3. **Diamond problem исчезает** — нет множественного наследования
4. **Код более понятен** — структура явно показывает, что содержит объект

```go
// Композиция проще
type Logger struct {
    level string
}

type Service struct {
    Logger  // Явно: Service содержит Logger
    data   map[string]interface{}
}

// Чем
type ServiceBase struct { ... }
class Service extends ServiceBase { ... }  // Скрытая зависимость
```

---

### Q: Как мне использовать ООП паттерны в Go?

**A:** Go поддерживает основные паттерны:

- **Factory:** функции `NewXxx`
- **Singleton:** одна переменная на уровне пакета
- **Observer:** интерфейсы + слайсы реципиентов
- **Strategy:** интерфейсы как стратегии
- **Decorator:** embedding типов
- **Repository:** интерфейсы для доступа к данным
- **Middleware:** функции высшего порядка

```go
// Strategy pattern
type SortStrategy interface {
    Sort([]int)
}

type Sorter struct {
    strategy SortStrategy
}

func (s *Sorter) Sort(arr []int) {
    s.strategy.Sort(arr)
}
```

---

### Q: Какие инструменты помогут писать хороший ООП код?

**A:**

```bash
# Проверка стиля
golangci-lint run ./...
go vet ./...

# Форматирование
gofmt -w ./...
goimports -w ./...

# Документация
godoc -http=:8080  # Откройте http://localhost:8080

# Тесты и покрытие
go test -v -cover ./...
go test -coverprofile=cover.out ./...
go tool cover -html=cover.out

# Профилирование
go build -cpuprofile=cpu.prof ./cmd/myapp
go tool pprof cpu.prof

# Статический анализ памяти
go install golang.org/x/exp/cmd/structlayout@latest
structlayout ./mypackage
```

---

**Спасибо за вопросы! Если появились новые, обращайтесь к документации: https://go.dev/doc/**
