# Swagger / OpenAPI: REST API с автодокументацией

Учебный пример REST API пользователей с **code-first** документацией: аннотации swag в коде → генерация OpenAPI → Swagger UI.

## Что демонстрирует

- описание API комментариями `@Summary`, `@Router`, `@Param`, …;
- генерацию `docs/` инструментом `swag init`;
- подключение UI через `http-swagger` на префиксе `/swagger/`;
- тот же предметный домен, что в `../4_rest` (users + JSON), плюс контракт для людей и инструментов.

## Структура

```text
swaggo_example/
├── README.md
├── main.go
├── go.mod
├── go.sum
└── docs/
    ├── docs.go
    ├── swagger.json
    └── swagger.yaml
```

Импорт `_ ".../docs"` нужен ради side-effect: зарегистрировать сгенерированную спецификацию для UI.

## API эндпоинты

| Метод | Путь | Описание |
|-------|------|----------|
| `GET` | `/users` | Список всех пользователей |
| `GET` | `/users/{id}` | Пользователь по ID |
| `POST` | `/users` | Создание пользователя, тело `{"name":"..."}` |

Интерактивная документация: `http://localhost:8080/swagger/index.html`

> В отличие от `../4_rest`, здесь пока нет `PUT`/`DELETE` в коде и в спецификации (их можно добавить по заданию корневого README). Маршрутизация здесь упрощённая (ветвление по `r.Method` и разбор пути), без шаблонов Go 1.22.

## Установка инструментов

Зависимости модуля:

```bash
go mod tidy
```

CLI для генерации документации (один раз в окружении разработчика):

```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

Убедитесь, что `$(go env GOPATH)/bin` есть в `PATH`.

## Генерация документации

После изменения аннотаций в `main.go`:

```bash
swag init -g main.go -o ./docs
```

- `-g main.go` - файл с `@title` / `@host` и точкой входа;
- `-o ./docs` - каталог артефактов OpenAPI.

Сгенерированные файлы **имеет смысл хранить в репозитории** и обновлять в том же PR, что меняет API: иначе Swagger UI врёт относительно кода.

## Запуск

```bash
go run main.go
```

- REST API: http://localhost:8080  
- Swagger UI: http://localhost:8080/swagger/index.html  

В UI можно выполнять запросы прямо из браузера («Try it out»).

## Примеры curl

```bash
curl -s http://localhost:8080/users
curl -s http://localhost:8080/users/1
curl -s -X POST -H "Content-Type: application/json" \
  -d '{"name":"Charlie"}' http://localhost:8080/users
```

## Разбор аннотаций

Над `main`:

```go
// @title Users API
// @version 1.0
// @host localhost:8080
// @BasePath /
```

Над handler:

```go
// @Summary Get user by ID
// @Param id path int true "User ID"
// @Success 200 {object} User
// @Failure 404 {string} string "User not found"
// @Router /users/{id} [get]
```

Чем полнее описаны схемы ответов, обязательные поля и ошибки, тем полезнее контракт для клиентов и тестов.

## Ограничения учебного примера

- нет PUT/DELETE и неполный CRUD относительно `../4_rest`;
- разбор `/users/{id}` через срез пути хрупче, чем `PathValue`;
- нет схем безопасности (Basic/JWT) в OpenAPI;
- in-memory store, данные сбрасываются при рестарте.

## Связь с лекцией

Тема OpenAPI / Swagger (лекция 2). Логичная связка: сначала понять REST на `../4_rest`, затем зафиксировать контракт здесь.
