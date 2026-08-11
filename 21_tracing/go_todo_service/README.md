# TODO-сервис с распределённой трассировкой

Учебный Go-сервис демонстрирует REST, gRPC и распределённую трассировку через OpenTelemetry Collector и Jaeger.

## Архитектура

```text
HTTP/gRPC client
       │
       ▼
Go Todo Service ── OTLP/HTTP ──► OpenTelemetry Collector ── OTLP/gRPC ──► Jaeger
       │
       ▼
   PostgreSQL
```

HTTP instrumentation создаётся `otelgin`, gRPC instrumentation — `otelgrpc`. Прикладные слои service и repository создают ручные дочерние spans. Логи HTTP-запросов содержат `trace_id`, `span_id` и `request_id`.

## Быстрый запуск

```bash
docker compose up --build -d
docker compose ps
curl -H 'X-User-ID: 1' http://localhost:8080/api/v1/todos
```

Адреса:

- REST API: <http://localhost:8080>;
- gRPC: `localhost:50051`;
- Jaeger UI: <http://localhost:16686>;
- Collector health: <http://localhost:13133>.

Остановка с сохранением данных:

```bash
docker compose down
```

Удаление контейнеров вместе с данными PostgreSQL:

```bash
docker compose down -v
```

## REST API

Все TODO-запросы требуют заголовок `X-User-ID`.

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H 'Content-Type: application/json' \
  -H 'X-User-ID: 1' \
  -d '{"title":"Изучить tracing","description":"Открыть Jaeger","done":false}'

curl -H 'X-User-ID: 1' http://localhost:8080/api/v1/todos
```

Проверки состояния:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/ready
```

`/health` и `/ready` намеренно исключены из трассировки.

## gRPC API

gRPC reflection включён в Compose.

```bash
grpcurl -plaintext localhost:50051 list todo.TodoService
grpcurl -plaintext -H 'user-id: 1' localhost:50051 todo.TodoService/ListTodos
grpcurl -plaintext -H 'user-id: 1' localhost:50051 todo.TodoService/ListTodosStream
```

Сервис содержит unary RPC, серверный поток и клиентский поток `BulkCreateTodos`.

## Конфигурация трассировки

| Переменная | Назначение | Значение для запуска на хосте |
|---|---|---|
| `TRACE_ENABLED` | Включить экспорт | `true` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP/HTTP endpoint | `http://localhost:4318` |
| `OTEL_EXPORTER_OTLP_INSECURE` | Разрешить соединение без TLS | `true` |
| `OTEL_SERVICE_NAME` | Имя сервиса в Jaeger | `go-todo-service` |
| `SERVICE_VERSION` | Версия сервиса | `1.0.0` |
| `DEPLOYMENT_ENVIRONMENT` | Окружение | `development` |
| `TRACE_SAMPLE_RATIO` | Доля корневых traces, от 0 до 1 | `1` |

В Compose приложение обращается к `http://otel-collector:4318`, поскольку `localhost` внутри контейнера означал бы сам контейнер приложения.

## Проверки

```bash
go test ./...
go test -race ./...
go vet ./...
docker compose config --quiet
```

Диагностика контейнеров:

```bash
docker compose ps -a
docker compose logs app otel-collector jaeger postgres
```

Подробная теория и практическое задание находятся в [лекции](../lecture.md).
