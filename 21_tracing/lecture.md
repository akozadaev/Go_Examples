# Распределённая трассировка в Go

## 1. Зачем нужна трассировка

Логи показывают сообщения отдельного процесса, метрики — состояние системы в целом, а трассировка — путь конкретного запроса и время каждой операции.

- **trace** — полный путь запроса;
- **span** — отдельная операция;
- **trace ID** — идентификатор всего пути;
- **span ID** — идентификатор одной операции;
- **parent/child** — связь операций;
- **attribute** — индексируемое свойство span;
- **event** — событие внутри span;
- **status** — результат операции: `Unset`, `Ok` или `Error`.

## 2. Архитектура примера

OpenTelemetry — нейтральный стандарт и набор библиотек. Приложение создаёт spans через OTel SDK и отправляет их по OTLP, не связывая прикладной код напрямую с Jaeger.

```text
HTTP или gRPC клиент
        │ W3C Trace Context
        ▼
Go Todo Service ── OTLP/HTTP :4318 ──► OTel Collector
        │                                  │ batch
        ▼                                  ▼ OTLP/gRPC :4317
    PostgreSQL                           Jaeger ──► UI :16686
```

- `app` создаёт spans;
- `postgres` хранит TODO;
- `otel-collector` принимает, обрабатывает и пересылает телеметрию;
- `jaeger` хранит traces в памяти и предоставляет UI.

In-memory traces Jaeger исчезают при пересоздании контейнера.

## 3. Контекст и propagation

Активный span хранится в `context.Context`. Контекст передаётся через handler → service → repository. Замена его на `context.Background()` в середине запроса разорвёт trace.

Между процессами используется W3C Trace Context:

```text
traceparent: 00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01
```

Серверная instrumentation извлекает входной контекст. Для продолжения trace при исходящем HTTP или gRPC-вызове нужна соответствующая клиентская instrumentation.

## 4. Инициализация SDK

`pkg/trace/trace_routes.go`:

1. создаёт OTLP/HTTP exporter;
2. формирует resource с именем, версией и окружением сервиса;
3. настраивает sampler и batch exporter;
4. устанавливает `TracerProvider` и W3C propagator;
5. при завершении вызывает `Shutdown` для отправки накопленных spans.

При `TRACE_ENABLED=false` используется безопасный no-op provider.

| Переменная | Назначение | Запуск Go на хосте | Compose |
|---|---|---|---|
| `TRACE_ENABLED` | экспорт traces | `true` | `true` |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP/HTTP endpoint | `http://localhost:4318` | `http://otel-collector:4318` |
| `OTEL_EXPORTER_OTLP_INSECURE` | соединение без TLS | `true` | `true` |
| `OTEL_SERVICE_NAME` | имя в Jaeger | `go-todo-service` | `go-todo-service` |
| `SERVICE_VERSION` | версия сервиса | `1.0.0` | `1.0.0` |
| `DEPLOYMENT_ENVIRONMENT` | окружение | `development` | `development` |
| `TRACE_SAMPLE_RATIO` | доля корневых traces | `1` | `1` |

## 5. Автоматическая instrumentation

### Gin

`otelgin.Middleware` создаёт серверный HTTP span, извлекает входной контекст и записывает стандартные HTTP-атрибуты. Middleware выполняется перед прикладным логгером, поэтому логи получают `trace_id` и `span_id`. `/health` и `/ready` исключены, чтобы probes не создавали шум.

### gRPC

`grpc.StatsHandler(otelgrpc.NewServerHandler())` работает с unary и streaming RPC. Один серверный span охватывает полный вызов, включая время жизни потока.

### Repository

Repository создаёт ручные spans вокруг операций. Это показывает границы слоя, но не детализацию SQL-драйвера. Для production можно добавить совместимую instrumentation `database/sql`, `pgx` или GORM, предварительно проверив зависимости и политику записи SQL.

## 6. Ручные spans

Фактический стиль проекта:

```go
ctx, span := tracer.Start(ctx, "service.CreateTodo")
defer span.End()

result, err := repository.Create(ctx, todo)
if err != nil {
    span.RecordError(err)
    span.SetStatus(codes.Error, "создание задачи")
    return nil, err
}
return result, nil
```

Instrumentation scope задаётся именем Go-пакета, а `service.name` — именем развёрнутого приложения. Успешный span обычно не требует явного `codes.Ok`: `Unset` означает отсутствие известной ошибки.

## 7. Ошибки, безопасность и кардинальность

`RecordError(err)` добавляет событие, но не меняет status. Для настоящего сбоя обычно также вызывается `SetStatus(codes.Error, ...)`. Ошибки оборачивают через `%w`, чтобы сохранить причину.

В spans нельзя без необходимости помещать:

- тела запросов;
- пароли, токены, cookie и персональные данные;
- заголовки и описания TODO;
- URL с секретами;
- произвольные неограниченные тексты.

`user.id` и `todo.id` высококардинальны. В примере они оставлены для обучения; в production это требует оценки безопасности и стоимости.

## 8. Sampling

`ParentBased(TraceIDRatioBased(...))` соблюдает решение родителя и применяет вероятность к новому корневому trace. Значение `1` удобно для занятия, но дорого под нагрузкой.

Head sampling принимает решение в начале запроса. Tail sampling в Collector может сохранить ошибочные или медленные traces после завершения, но требует больше памяти и более сложной маршрутизации.

## 9. Практическая проверка

```bash
cd /home/akozadaev/projects/akozadaev/go/ibs/Go_Examples/21_tracing/go_todo_service
docker compose up --build -d
docker compose ps
curl http://localhost:8080/ready
curl -H 'X-User-ID: 1' http://localhost:8080/api/v1/todos
```

Создание задачи:

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H 'Content-Type: application/json' \
  -H 'X-User-ID: 1' \
  -d '{"title":"Изучить tracing","description":"Открыть Jaeger","done":false}'
```

Откройте <http://localhost:16686>, выберите `go-todo-service` и нажмите **Find Traces**. Из-за batch exporter trace может появиться через несколько секунд.

Проверка gRPC:

```bash
grpcurl -plaintext -H 'user-id: 1' localhost:50051 todo.TodoService/ListTodos
grpcurl -plaintext -H 'user-id: 1' localhost:50051 todo.TodoService/ListTodosStream
```

Ожидаемое дерево REST trace:

```text
POST /api/v1/todos
└── service.CreateTodo
    └── repository.CreateTodo
```

## 10. Диагностика

Если traces не видны:

1. отправьте запрос не на `/health` и не на `/ready`;
2. проверьте `TRACE_ENABLED` и endpoint;
3. выполните `docker compose ps`;
4. изучите `docker compose logs app otel-collector jaeger`;
5. подождите отправки batch;
6. проверьте имя сервиса в Jaeger;
7. завершайте приложение через SIGTERM, чтобы SDK выполнил `Shutdown`.

## 11. Production checklist

- TLS и аутентификация между SDK, Collector и backend;
- постоянное хранилище вместо in-memory;
- sampling и лимиты объёма;
- несколько Collector, очередь и retry;
- фильтрация секретов и персональных данных;
- единые semantic conventions;
- связь логов по `trace_id` и `span_id`;
- мониторинг Collector;
- закреплённые версии образов и регулярное обновление;
- нагрузочные тесты стоимости instrumentation.

## 12. Практическое задание

1. Создайте TODO через HTTP и объясните дерево spans.
2. Повторите запрос через gRPC и сравните корневые spans.
3. Передайте собственный `traceparent` и проверьте продолжение trace.
4. Установите `TRACE_SAMPLE_RATIO=0` и объясните результат.
5. Добавьте задержку в repository, найдите её по timeline, затем удалите.

## Контрольные вопросы

1. Чем trace отличается от лога и метрики?
2. Почему нельзя создавать новый контекст в середине запроса?
3. За что отвечают SDK, Collector и Jaeger?
4. Почему тело запроса — плохой атрибут span?
5. Чем head sampling отличается от tail sampling?
6. Зачем вызывать `TracerProvider.Shutdown`?
