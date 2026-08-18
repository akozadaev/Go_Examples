# Кодогенерация с `easyjson`

## Назначение примера

Универсальный `encoding/json` должен работать с заранее неизвестными типами и поэтому использует динамические механизмы. В горячем пути высоконагруженного приложения стоимость сериализации и аллокаций иногда становится заметной.

Этот пример показывает применение [`easyjson`](https://pkg.go.dev/github.com/mailru/easyjson) к структурам `User` и `Profile`.

- **Решаемая проблема:** уменьшить динамическую работу и количество аллокаций при частой JSON-сериализации и десериализации известных структур.
- **Что генерируется:** специализированные encoder/decoder и методы `MarshalJSON`, `UnmarshalJSON`, `MarshalEasyJSON`, `UnmarshalEasyJSON`.
- **Когда полезно:** после измерения профилем и benchmark в сервисах, где обработка JSON действительно является существенной частью нагрузки.

Generated-код обращается к полям напрямую, не исследуя структуру через runtime-reflection на каждом вызове. Это не означает автоматического ускорения любого приложения: необходимо проверить совместимость поведения и измерить конкретные данные.

## Файлы

| Роль | Файл относительно корня `go_code_generation` |
|---|---|
| Структуры `User`, `Profile` и директива | `easyjson_example/user.go` |
| Сгенерированный код | `easyjson_example/user_easyjson.go` |
| Демонстрация | `cmd/easyjson/main.go` |
| Тест корректности и benchmark | `easyjson_example/benchmark_test.go` |
| Дополнительное описание | `easyjson_example/benchmark.md` |

Директива в `user.go`:

```go
//go:generate easyjson -all user.go
```

`-all` просит обработать все структуры файла — в данном случае `User` и `Profile`.

## Установка и генерация

Runtime-зависимость указана в `go.mod`, но CLI устанавливается отдельно:

```bash
go install github.com/mailru/easyjson/easyjson@latest
go generate ./easyjson_example
```

В воспроизводимом проекте версию фиксируют и согласуют с версией библиотеки в `go.mod`. Из `easyjson_example` можно вызвать `go generate`. Посмотреть команду без выполнения:

```bash
go generate -n ./easyjson_example
```

Фактически запускается `easyjson -all user.go`, а результат попадает рядом с исходником:

```text
easyjson_example/user_easyjson.go
```

## Что генерируется

Файл содержит специализированные encoder/decoder и методы для обеих структур:

```go
func (v User) MarshalJSON() ([]byte, error)
func (v *User) UnmarshalJSON(data []byte) error
func (v User) MarshalEasyJSON(w *jwriter.Writer)
func (v *User) UnmarshalEasyJSON(l *jlexer.Lexer)
```

Из-за методов `MarshalJSON`/`UnmarshalJSON` пакет `encoding/json` тоже может делегировать работу easyjson. Поэтому benchmark использует отдельные `plainUser` и `plainProfile` без generated-методов для честного стандартного baseline.

## Проверка

```bash
go run ./cmd/easyjson
go test ./easyjson_example
go test -bench=. -benchmem ./easyjson_example
```

`TestJSONImplementationsAreEquivalent` проверяет семантическое равенство JSON. Benchmark измеряет `ns/op`, `B/op`, `allocs/op`; конкретное ускорение зависит от данных и платформы.

`user_easyjson.go` полностью перезаписывается генератором и исключён из Git правилом `*_easyjson.go`. Менять следует структуры, теги или параметры в `user.go`, но не generated-файл.
