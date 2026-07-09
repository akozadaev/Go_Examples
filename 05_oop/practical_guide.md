# Практическое руководство: ООП в Golang

**Назначение:** Дополнительные примеры, рецепты и советы для применения ООП в production коде.

---

## Содержание

1. [Разработка типов](#разработка-типов)
2. [Интерфейсы в production](#интерфейсы-в-production)
3. [Контроль доступа](#контроль-доступа)
4. [Обработка ошибок](#обработка-ошибок)
5. [Паттерны проектирования](#паттерны-проектирования)
6. [Оптимизация памяти](#оптимизация-памяти)
7. [Тестирование](#тестирование)
8. [Чек-лист для code review](#чек-лист-для-code-review)

---

## Разработка типов

### Правило: Функции вместо методов

**Когда использовать функции:**

```go
// Хорошо: свободная функция
func Distance(p1, p2 Point) float64 {
    dx := p2.X - p1.X
    dy := p2.Y - p1.Y
    return math.Sqrt(dx*dx + dy*dy)
}

// Используется
dist := Distance(p1, p2)
```

**Когда использовать методы:**

```go
// Хорошо: метод если операция связана с типом
type Point struct {
    X, Y float64
}

func (p Point) Length() float64 {
    return math.Sqrt(p.X*p.X + p.Y*p.Y)
}

// Используется
dist := p.Length()
```

**Правило:** если первый аргумент функции совпадает по смыслу с receiver'ом, используйте метод.

---

### Конструкторы

**Паттерн 1: Простой конструктор**

```go
type User struct {
    id    int64
    name  string
    email string
}

func NewUser(name, email string) *User {
    return &User{
        id:    generateID(),
        name:  name,
        email: email,
    }
}
```

**Паттерн 2: Конструктор с опциями (Functional Options)**

```go
type Logger struct {
    level  string
    output io.Writer
    prefix string
}

type LoggerOption func(*Logger)

func WithLevel(level string) LoggerOption {
    return func(l *Logger) {
        l.level = level
    }
}

func WithPrefix(prefix string) LoggerOption {
    return func(l *Logger) {
        l.prefix = prefix
    }
}

func NewLogger(opts ...LoggerOption) *Logger {
    l := &Logger{
        level:  "info",
        output: os.Stderr,
        prefix: "",
    }
    for _, opt := range opts {
        opt(l)
    }
    return l
}

// Использование
logger := NewLogger(
    WithLevel("debug"),
    WithPrefix("[APP] "),
)
```

**Паттерн 3: Builder для сложных объектов**

```go
type Config struct {
    Host     string
    Port     int
    Database string
    Username string
    Password string
    Timeout  time.Duration
}

type ConfigBuilder struct {
    host     string
    port     int
    database string
    username string
    password string
    timeout  time.Duration
}

func NewConfigBuilder() *ConfigBuilder {
    return &ConfigBuilder{
        host:     "localhost",
        port:     5432,
        timeout:  30 * time.Second,
    }
}

func (b *ConfigBuilder) Host(h string) *ConfigBuilder {
    b.host = h
    return b
}

func (b *ConfigBuilder) Port(p int) *ConfigBuilder {
    b.port = p
    return b
}

func (b *ConfigBuilder) Build() *Config {
    return &Config{
        Host:     b.host,
        Port:     b.port,
        Database: b.database,
        Username: b.username,
        Password: b.password,
        Timeout:  b.timeout,
    }
}

// Использование
config := NewConfigBuilder().
    Host("example.com").
    Port(3306).
    Build()
```

---

## Интерфейсы в production

### Принцип: Small, focused interfaces

```go
// Хорошо: маленкие, сфокусированные интерфейсы
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

// Плохо: огромный интерфейс (God Interface)
type Database interface {
    Query(string) ([]Row, error)
    Exec(string) error
    BeginTx(context.Context) (Tx, error)
    Close() error
    Ping() error
    Stats() Stats
    // ... ещё 20 методов
}
```

### Принцип: Define interfaces, not implementations

```go
// Неправильно: экспортируем конкретный тип
type PostgresDB struct { ... }
func NewPostgresDB() *PostgresDB { ... }

// Правильно: экспортируем интерфейс
type Database interface {
    Query(ctx context.Context, q string) ([]Row, error)
    Exec(ctx context.Context, q string) error
}

type postgresDB struct { ... }
func NewDatabase() Database {
    return &postgresDB{...}
}
```

### Принцип: Accept interfaces, return concrete types

```go
// Правильно: принимаем интерфейс (гибко)
func Process(r io.Reader) ([]byte, error) {
    return ioutil.ReadAll(r)
}

// Правильно: возвращаем конкретный тип (точно)
func NewFileReader(path string) (*os.File, error) {
    return os.Open(path)
}

// Неправильно: возвращаем интерфейс
func Process() io.Reader {
    return os.Stdin  // Непредсказуемо
}
```

---

## Контроль доступа

### Защита инвариантов через setter'ы

```go
type BankAccount struct {
    balance     float64
    owner       string
    opened      time.Time
    // lastAccess не нужен снаружи
}

// Getter для публичного поля
func (a *BankAccount) Balance() float64 {
    return a.balance
}

// Setter с валидацией инварианта
func (a *BankAccount) Deposit(amount float64) error {
    if amount <= 0 {
        return fmt.Errorf("amount must be positive, got %v", amount)
    }
    a.balance += amount
    return nil
}

// Getter для приватных данных
func (a *BankAccount) Owner() string {
    return a.owner
}

// Операция требует обновления связанных полей
func (a *BankAccount) SetOwner(newOwner string) error {
    if newOwner == "" {
        return errors.New("owner cannot be empty")
    }
    a.owner = newOwner
    return nil
}
```

### Инкапсуляция состояния

```go
// Неправильно: состояние экспортировано
type Order struct {
    Status string  // Может быть любое значение!
}

// Правильно: состояние контролируется через методы
type Order struct {
    status OrderStatus  // Приватное поле с типом
}

type OrderStatus int

const (
    OrderPending   OrderStatus = iota
    OrderProcessing
    OrderCompleted
    OrderCancelled
)

func (o *Order) Status() OrderStatus {
    return o.status
}

func (o *Order) Complete() error {
    if o.status != OrderProcessing {
        return fmt.Errorf("cannot complete order in %v status", o.status)
    }
    o.status = OrderCompleted
    return nil
}
```

---

## Обработка ошибок

### Определение кастомных ошибок

```go
// Ошибка 1: значение ошибки (Sentinel Error)
var (
    ErrNotFound      = errors.New("resource not found")
    ErrUnauthorized  = errors.New("access denied")
    ErrInvalidInput  = errors.New("invalid input provided")
)

// Ошибка 2: тип ошибки (Error Type)
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error in field %q: %s", e.Field, e.Message)
}

// Ошибка 3: поведенческая ошибка (Error Behavior)
type TemporaryError interface {
    error
    Temporary() bool
}

type retryableError struct {
    msg       string
    retryable bool
}

func (e *retryableError) Error() string { return e.msg }
func (e *retryableError) Temporary() bool { return e.retryable }

// Использование
func SomeOperation() error {
    if someValidationFailed {
        return &ValidationError{
            Field:   "email",
            Message: "invalid email format",
        }
    }
    return nil
}

// Клиентский код
if err := SomeOperation(); err != nil {
    var ve *ValidationError
    if errors.As(err, &ve) {
        fmt.Printf("Validation failed on %s\n", ve.Field)
    }
}
```

### Оборачивание ошибок

```go
func ReadConfig(path string) (*Config, error) {
    data, err := os.ReadFile(path)
    if err != nil {
        // Оборачиваем ошибку контекстом
        return nil, fmt.Errorf("failed to read config file at %s: %w", path, err)
    }
    
    var cfg Config
    if err := json.Unmarshal(data, &cfg); err != nil {
        return nil, fmt.Errorf("failed to parse config: %w", err)
    }
    
    return &cfg, nil
}

// Клиентский код
config, err := ReadConfig("config.json")
if err != nil {
    // Можем распаковать оригинальную ошибку
    if os.IsNotExist(err) {
        fmt.Println("Config file not found")
    }
}
```

---

## Паттерны проектирования

### Repository Pattern

```go
// Интерфейс репозитория
type UserRepository interface {
    GetByID(ctx context.Context, id int64) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Save(ctx context.Context, u *User) error
    Delete(ctx context.Context, id int64) error
}

// Реализация (приватная)
type postgresUserRepository struct {
    db *sql.DB
}

func (r *postgresUserRepository) GetByID(ctx context.Context, id int64) (*User, error) {
    query := "SELECT id, name, email FROM users WHERE id = $1"
    // ... реализация
}

// Конструктор
func NewUserRepository(db *sql.DB) UserRepository {
    return &postgresUserRepository{db: db}
}

// Service использует интерфейс
type UserService struct {
    repo UserRepository
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*User, error) {
    user, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("failed to get user: %w", err)
    }
    return user, nil
}
```

### Middleware Pattern

```go
type Handler func(context.Context, *Request) (*Response, error)

// Middleware функция
func LoggingMiddleware(h Handler) Handler {
    return func(ctx context.Context, req *Request) (*Response, error) {
        log.Printf("Request: %v", req)
        resp, err := h(ctx, req)
        if err != nil {
            log.Printf("Error: %v", err)
        }
        return resp, err
    }
}

func AuthMiddleware(h Handler) Handler {
    return func(ctx context.Context, req *Request) (*Response, error) {
        if req.Token == "" {
            return nil, fmt.Errorf("missing authentication token")
        }
        return h(ctx, req)
    }
}

// Использование
handler := LoggingMiddleware(AuthMiddleware(myHandler))
resp, err := handler(ctx, req)
```

### Observer Pattern

```go
type Observer interface {
    Update(event Event)
}

type EventPublisher struct {
    observers []Observer
    mu        sync.RWMutex
}

func (p *EventPublisher) Subscribe(o Observer) {
    p.mu.Lock()
    p.observers = append(p.observers, o)
    p.mu.Unlock()
}

func (p *EventPublisher) Publish(event Event) {
    p.mu.RLock()
    observers := make([]Observer, len(p.observers))
    copy(observers, p.observers)
    p.mu.RUnlock()
    
    for _, obs := range observers {
        go obs.Update(event)  // Асинхронно уведомляем
    }
}

// Конкретный observer
type EmailNotifier struct {
    email string
}

func (n *EmailNotifier) Update(event Event) {
    log.Printf("Sending email to %s about %v", n.email, event)
}

// Использование
pub := &EventPublisher{}
pub.Subscribe(&EmailNotifier{email: "admin@example.com"})
pub.Publish(Event{Type: "UserCreated"})
```

---

## Оптимизация памяти

### Анализ выравнивания

```go
// Инструмент для анализа
// go install golang.org/x/exp/cmd/structlayout@latest
// structlayout ./mypackage

type Example1 struct {
    flag1  bool      // 1 byte
    value1 int64     // 8 bytes
    flag2  bool      // 1 byte
}
// Size: 24 bytes (много padding)

type Example2 struct {
    value1 int64     // 8 bytes
    value2 int32     // 4 bytes
    flag1  bool      // 1 byte
    flag2  bool      // 1 byte
}
// Size: 16 bytes (меньше padding)

// Правило: от большего к меньшему
type Optimized struct {
    value1 int64
    value2 int32
    value3 int16
    flag1  bool
    flag2  bool
}
```

### Пулы объектов для производительности

```go
import "sync"

type RequestPool struct {
    pool *sync.Pool
}

func NewRequestPool() *RequestPool {
    return &RequestPool{
        pool: &sync.Pool{
            New: func() interface{} {
                return &Request{}
            },
        },
    }
}

func (p *RequestPool) Get() *Request {
    return p.pool.Get().(*Request)
}

func (p *RequestPool) Put(r *Request) {
    r.Reset()
    p.pool.Put(r)
}

// Использование
pool := NewRequestPool()

req := pool.Get()
defer pool.Put(req)

// ... используем req
```

---

## Тестирование

### Mock интерфейсы

```go
// Реальный интерфейс
type Database interface {
    Query(ctx context.Context, q string) ([]Row, error)
}

// Mock для тестирования
type MockDatabase struct {
    QueryFunc func(ctx context.Context, q string) ([]Row, error)
    QueryCalls int
}

func (m *MockDatabase) Query(ctx context.Context, q string) ([]Row, error) {
    m.QueryCalls++
    return m.QueryFunc(ctx, q)
}

// Тест
func TestGetUser(t *testing.T) {
    mockDB := &MockDatabase{
        QueryFunc: func(ctx context.Context, q string) ([]Row, error) {
            return []Row{{Name: "Alice"}}, nil
        },
    }
    
    service := &UserService{db: mockDB}
    user, err := service.GetUser(context.Background(), 1)
    
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if mockDB.QueryCalls != 1 {
        t.Errorf("expected 1 query call, got %d", mockDB.QueryCalls)
    }
}
```

### Table-driven тесты

```go
func TestDiscount(t *testing.T) {
    tests := []struct {
        name      string
        amount    float64
        percent   int
        expected  float64
        shouldErr bool
    }{
        {
            name:     "valid discount",
            amount:   100,
            percent:  10,
            expected: 90,
        },
        {
            name:      "negative amount",
            amount:    -100,
            percent:   10,
            shouldErr: true,
        },
        {
            name:      "percent too high",
            amount:    100,
            percent:   150,
            shouldErr: true,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result, err := CalculateDiscount(tt.amount, tt.percent)
            
            if tt.shouldErr && err == nil {
                t.Error("expected error, got nil")
            }
            if !tt.shouldErr && err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            if !tt.shouldErr && result != tt.expected {
                t.Errorf("expected %v, got %v", tt.expected, result)
            }
        })
    }
}
```

---

## Чек-лист для code review

### Структуры и типы

- [ ] Используется композиция вместо наследования?
- [ ] Поля правильно расположены (от большего к меньшему)?
- [ ] Определены конструкторы с валидацией?
- [ ] Используются неэкспортируемые поля для защиты?

### Методы и receivers

- [ ] Выбран правильный тип receiver (value vs pointer)?
- [ ] Все методы типа используют один стиль receiver'ов?
- [ ] Методы не изменяют receiver без необходимости?
- [ ] Используются правильные соглашения имен (Get, Set, Is)?

### Инкапсуляция

- [ ] API экспортирует интерфейсы, не реализации?
- [ ] Приватные поля защищены setter'ами?
- [ ] Getter'ы скрывают детали реализации?
- [ ] Инварианты объекта проверяются?

### Обработка ошибок

- [ ] Ошибки определены явно (Sentinel или Type)?
- [ ] Ошибки обернуты с контекстом (`%w`)?
- [ ] Корректно используются `errors.Is` и `errors.As`?
- [ ] Нет игнорирования ошибок (`_ = ...`)?

### Интерфейсы

- [ ] Интерфейсы маленькие и сфокусированные?
- [ ] Экспортируются интерфейсы, не конкретные типы?
- [ ] Методы принимают интерфейсы, возвращают типы?
- [ ] Интерфейсы определены близко к месту использования?

### Производительность

- [ ] Нет ненужного копирования структур?
- [ ] Большие структуры передаются по указателю?
- [ ] Не создаются циклические выделения памяти?
- [ ] Используются пулы объектов для горячих путей?

---

## Полезные команды

```bash
# Проверка выравнивания структур
go install golang.org/x/exp/cmd/structlayout@latest
structlayout ./mypackage

# Анализ кода
go vet ./...
golangci-lint run ./...

# Профилирование памяти
go build -memprofile=mem.prof
go tool pprof mem.prof

# Исправление форматирования
gofmt -w ./...
goimports -w ./...

# Запуск тестов с покрытием
go test -cover ./...
go test -coverprofile=cover.out ./...
go tool cover -html=cover.out
```

---

## Дополнительные ресурсы

### На русском
- "Язык программирования Go" (тур)
- "Отправляющийся в Golang" (Donovan, Kernighan)

### На английском
- "The Go Programming Language" (Donovan, Kernighan)
- "Concurrency in Go" (Katherine Cox-Buday)
- https://go.dev/doc/effective_go
- https://github.com/golang/go/wiki/CodeReviewComments

### Инструменты
- [Effective Go](https://go.dev/doc/effective_go)
- [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
