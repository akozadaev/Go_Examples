# Генератор мапперов для Go

Этот проект демонстрирует применение кодогенерации в Go для автоматического создания методов преобразования между моделями и DTO.

## Структура проекта

- `cmd/mapper/` - генератор мапперов
- `mapper_example/` - примеры использования генератора
- `cmd/mapper-demo/` - демонстрация работы генератора

## Использование

### 1. Определение структур

Создайте модели и DTO с соответствующими тегами:

```go
type User struct {
    ID        int64     `mapper:"id"`
    Name      string    `mapper:"name"`
    CreatedAt time.Time `mapper:"created_at"`
}

type UserDTO struct {
    ID        int64  `json:"id"`
    Name      string `json:"name"`
    CreatedAt string `json:"created_at"`
}
```

### 2. Добавление директивы go:generate

```go
//go:generate go run ../cmd/mapper/main.go -type=User -dto=UserDTO -output=user_mapper_gen.go user.go
```

### 3. Генерация кода

```bash
go generate ./mapper_example
```

или для всех пакетов:

```bash
make generate
```

### 4. Использование сгенерированных методов

```go
// Преобразование модели в DTO
user := &User{ID: 1, Name: "Иван"}
dto := user.ToUserDTO()

// Обратное преобразование
userFromDTO := UserDTOToUser(dto)

// Преобразование слайсов
users := []User{...}
dtos := SliceToUserDTO(users)
```

## Особенности

- Автоматическое преобразование `time.Time` ↔ `string` (RFC3339)
- Поддержка слайсов структур
- Типобезопасность на этапе компиляции
- Интеграция с `go generate`

## Запуск демонстрации

```bash
make run
```

или

```bash
go run ./cmd/mapper-demo
```

