# Примеры логирования на Go

В каталоге лежат отдельные программы (`main_*.go`), каждая из которых демонстрирует свой способ логирования. Все файлы объявляют `package main` и запускаются по одному: `go run <файл>.go`.

**Модуль:** `github.com/akozadaev/go_logging`  
**Версия Go:** см. директиву `go` в `go.mod` (для `log/slog` нужна поддержка пакета в вашей версии компилятора).

**Зависимости:** `go.uber.org/zap`, `github.com/rs/zerolog`; стандартные `log` и `log/slog` подтягиваются из SDK.

Установка зависимостей после клонирования:

```bash
go mod download
```

или

```bash
go mod tidy
```

---

## Сводная таблица примеров

| Файл | Библиотека | Назначение |
|------|------------|------------|
| `main_simple_log.go` | `log` | Базовый вывод в stderr: `Println`, `Printf` при ошибке |
| `main_log_to_file.go` | `log` | Перенаправление вывода логгера в файл `app.log` |
| `main_slog.go` | `log/slog` | Стандартный текстовый лог и отдельный логгер с JSON-handler |
| `main_zap_interface.go` | Zap | Свой интерфейс `Logger` и реализация на `zap.NewProduction()` |
| `main_zerolog_json.go` | Zerolog | Глобальный логгер пакета `log` — JSON в stdout |
| `main_zerolog_console.go` | Zerolog | `ConsoleWriter`: время, уровень, caller, человекочитаемый вывод |
| `main_zerolog_context.go` | Zerolog + `net/http` | HTTP-сервер: логгер в `context`, извлечение в хендлере |

---

## `main_simple_log.go` — стандартный `log`

- Сообщение через `log.Println`.
- Имитация ошибки в `someFunction` и запись через `log.Printf` с `%v`.

Запуск:

```bash
go run main_simple_log.go
```

В Makefile: `make run-log-simple`.

---

## `main_log_to_file.go` — `log` в файл

- Создаётся файл `app.log`, вывод логгера переключается через `log.SetOutput`.
- При ошибке создания файла — `log.Fatal`.

Запуск:

```bash
go run main_log_to_file.go
```

После запуска появится `app.log`. Удаление: `make clean` или `rm -f app.log`.

В Makefile: `make run-log-file`.

---

## `main_slog.go` — `log/slog`

- Первый вызов: дефолтный `slog` (текстовый формат) с парами ключ–значение.
- Второй блок: `slog.NewJSONHandler(os.Stdout)` — сообщения `Info` и `Error` в JSON.

Запуск:

```bash
go run main_slog.go
```

В Makefile: `make run-slog`.

---

## `main_zap_interface.go` — Zap и интерфейс

- Интерфейс `Logger` с методами `Info` и `Error`.
- Тип `ZapLogger` оборачивает `*zap.Logger`; конфигурация production.
- В `main` вызывается `Sync()` у внутреннего zap-логгера после работы.

Запуск:

```bash
go run main_zap_interface.go
```

В Makefile: `make run-zap`.

---

## `main_zerolog_json.go` — Zerolog (JSON)

- Используется глобальный логгер из `github.com/rs/zerolog/log`.
- Одно информационное сообщение в JSON на stdout.

Запуск:

```bash
go run main_zerolog_json.go
```

В Makefile: `make run-zerolog-json`.

---

## `main_zerolog_console.go` — Zerolog (консоль)

- Локальный логгер: `zerolog.ConsoleWriter` на `stderr`, формат времени RFC3339.
- Уровень `Trace`, в цепочке `With()` добавлены timestamp и caller.

Запуск:

```bash
go run main_zerolog_console.go
```

В Makefile: `make run-zerolog-console`.

---

## `main_zerolog_context.go` — Zerolog и контекст HTTP

- Поднимается сервер на `:8080`, маршрут `/delete`.
- Middleware создаёт логгер с полем `user_id`, кладёт его в `context` запроса (`WithContext`).
- Хендлер читает логгер через `zerolog.Ctx(r.Context())` и пишет лог с `doc_id`.

Запуск сервера:

```bash
go run main_zerolog_context.go
```

Проверка запросом:

```bash
curl http://127.0.0.1:8080/delete
```

В Makefile: `make run-zerolog-http`.

---

## Makefile

Команда `make` или `make help` выводит список целей. Для примеров используются цели `run-*`, см. таблицу выше и `Makefile` в корне репозитория.

---

## Полезные ссылки

- [Zap](https://github.com/uber-go/zap)
- [Zerolog](https://github.com/rs/zerolog)
- [Документация `log/slog`](https://pkg.go.dev/log/slog)
