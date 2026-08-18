# Кодогенерация с `stringer`

## Назначение примера

В Go перечисления обычно объявляют как именованный целочисленный тип с набором констант. Без дополнительного метода при выводе такого значения виден номер — например, `0`, — а не понятное имя `StatusPending`.

Этот пример показывает, как [`stringer`](https://pkg.go.dev/golang.org/x/tools/cmd/stringer) автоматически создаёт метод `String()` для типа `Status`.

- **Решаемая проблема:** ручной метод `String()` получается шаблонным и легко перестаёт соответствовать константам после их изменения.
- **Что генерируется:** типобезопасный метод `func (Status) String() string`, удовлетворяющий `fmt.Stringer`.
- **Когда полезно:** логирование, диагностика, CLI, сообщения об ошибках и любой читаемый вывод enum-значений.

`String()` обычно предназначен для представления человеку. Если строка становится частью JSON, базы данных или сетевого протокола, её стабильность нужно проектировать и тестировать отдельно.

## Исходники и результат

| Роль | Файл относительно корня `go_code_generation` |
|---|---|
| Тип, константы и директива | `stringer_example/status.go` |
| Сгенерированный файл | `stringer_example/status_string.go` |
| Демонстрационная программа | `cmd/stringer/main.go` |

В `status.go` объявлены `StatusPending`, `StatusInProgress`, `StatusCompleted`, `StatusCancelled` и директива:

```go
//go:generate stringer -type=Status
```

## Установка и генерация

`go generate` не устанавливает инструмент. Сначала `stringer` должен появиться в `PATH`:

```bash
go install golang.org/x/tools/cmd/stringer@latest
```

В production-проекте вместо `@latest` закрепляют конкретную версию. Затем из корня модуля выполните:

```bash
go generate ./stringer_example
```

Из каталога `stringer_example` эквивалентная команда выглядит как `go generate`. Посмотреть план без изменения файлов:

```bash
go generate -n ./stringer_example
```

Директива фактически запускает `stringer -type=Status`. По умолчанию результат записывается рядом с исходником в `stringer_example/status_string.go`.

## Что находится в generated-файле

Файл содержит обычный Go-метод:

```go
func (i Status) String() string
```

После генерации `fmt.Println(StatusPending)` выводит `StatusPending`, а не `0`. Файл начинается с маркера `Code generated ... DO NOT EDIT` и полностью перезаписывается при следующем запуске.

## Проверка

```bash
go run ./cmd/stringer
go test ./stringer_example
```

Правило `.gitignore` исключает `*_string.go`, поэтому после чистого клонирования generated-файл нужно создать до сборки демонстрации.
