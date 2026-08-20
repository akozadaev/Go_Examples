# Cценарий для поиска утечки/деградации.
 Предполагаемая проблема - отсутствие пагинации в getAll
## 1. Запустить стек

Из корня проекта:

```bash
docker compose up --build -d
docker compose ps
```

Проверить доступность:

```bash
curl http://localhost:8080/health
curl http://localhost:8080/ready
curl http://localhost:9090/metrics
```

Интерфейсы:

- API: http://localhost:8080
- gRPC: `localhost:50051`
- metrics: http://localhost:9090/metrics
- pprof: http://localhost:9090/debug/pprof/
- Prometheus: http://localhost:9091
- Jaeger: http://localhost:16686

Посмотреть стартовые ошибки:

```bash
docker compose logs --tail=200 app
docker compose logs --tail=100 prometheus
docker compose logs --tail=100 otel-collector
```

## 2. Проверить обычный сценарий

```bash
make test-curl
make test-grpc
```

Либо минимально:

```bash
curl -H 'X-User-ID: 1' \
  http://localhost:8080/api/v1/todos
```

Создание записи:

```bash
curl -X POST http://localhost:8080/api/v1/todos \
  -H 'Content-Type: application/json' \
  -H 'X-User-ID: 1' \
  -d '{"title":"observation-test","description":"test","done":false}'
```

## 3. Убедиться, что метрики собираются

```bash
curl -fsS http://localhost:9090/metrics |
  rg 'todo_|go_goroutines|go_memstats|process_resident'
```

Основные метрики:

```text
todo_http_requests_total
todo_http_request_duration_seconds
todo_http_requests_in_flight
todo_grpc_requests_total
todo_grpc_request_duration_seconds
todo_grpc_requests_in_flight
go_goroutines
go_memstats_heap_inuse_bytes
go_memstats_heap_objects
process_resident_memory_bytes
```

Названия доступных DB-метрик:

```bash
curl -fsS http://localhost:9090/metrics | rg 'todo_db_'
```

В Prometheus можно выполнить запросы:

```promql
rate(todo_http_requests_total[1m])
```

```promql
histogram_quantile(
  0.95,
  sum by (le, route) (
    rate(todo_http_request_duration_seconds_bucket[5m])
  )
)
```

```promql
go_memstats_heap_inuse_bytes
```

```promql
process_resident_memory_bytes
```

```promql
go_goroutines
```

```promql
rate(go_gc_duration_seconds_count[5m])
```

## 4. Проверить корреляцию логов и traces

Отправить запрос с известным request ID:

```bash
curl -v \
  -H 'X-User-ID: 1' \
  http://localhost:8080/api/v1/todos
```

Посмотреть логи:

```bash
docker compose logs --since=2m app
```

У HTTP-, gRPC- и SQL-событий одного запроса должны совпадать:

```text
trace_id
```

SQL-лог дополнительно должен содержать:

```text
component=gorm
duration
rows
sql
```

Скопировать `trace_id` из лога и найти его в Jaeger:

1. Открыть http://localhost:16686.
2. Выбрать `go-todo-service`.
3. Выполнить поиск по Trace ID.
4. Проверить последовательность HTTP/gRPC → service → repository.

Для latency histogram экспортируются exemplars с `trace_id`. Их можно увидеть через Prometheus при использовании OpenMetrics.

## 5. Снять baseline

Создать каталог вне Git-репозитория:

```bash
mkdir -p /tmp/go-todo-profiles
```

Снять heap после принудительного GC:

```bash
curl -fsS \
  'http://localhost:9090/debug/pprof/heap?gc=1' \
  -o /tmp/go-todo-profiles/heap-baseline.pprof
```

Goroutine dump:

```bash
curl -fsS \
  'http://localhost:9090/debug/pprof/goroutine?debug=2' \
  -o /tmp/go-todo-profiles/goroutine-baseline.txt
```

Записать исходные показатели:

```bash
curl -fsS http://localhost:9090/metrics |
  rg 'go_goroutines|go_memstats_heap_inuse_bytes|go_memstats_heap_objects|process_resident_memory_bytes'
```

Если установлен `PPROF_TOKEN`, добавить:

```bash
-H "Authorization: Bearer $PPROF_TOKEN"
```

## 6. Провести обычный нагрузочный тест

Для начала использовать небольшую нагрузку:

```bash
PROFILE_REQUESTS=5000 \
PROFILE_CONCURRENCY=50 \
PROFILE_OUTPUT_DIR=/tmp/go-todo-profiles \
make profile-load
```

Одновременно наблюдать метрики в Prometheus:

```promql
rate(todo_http_requests_total[1m])
```

```promql
histogram_quantile(
  0.99,
  sum by (le) (
    rate(todo_http_request_duration_seconds_bucket[1m])
  )
)
```

```promql
go_goroutines
```

```promql
process_resident_memory_bytes
```

Обращать внимание на:

- монотонный рост heap после завершения нагрузки;
- goroutine, которые не возвращаются к baseline;
- увеличение DB connection wait;
- рост p95/p99;
- HTTP 5xx и gRPC-коды ошибок;
- длительные repository/SQL spans.

## 7. Проверить проблему неограниченного `GetAll`

Создать достаточное количество TODO:

```bash
seq 1 10000 | xargs -P 20 -I '{}' \
  curl -fsS -o /dev/null \
  -X POST http://localhost:8080/api/v1/todos \
  -H 'Content-Type: application/json' \
  -H 'X-User-ID: 1' \
  -d '{"title":"load-{}","description":"memory test payload","done":false}'
```

Затем снять CPU profile прямо во время большого `GetAll`:

```bash
curl -fsS \
  'http://localhost:9090/debug/pprof/profile?seconds=30' \
  -o /tmp/go-todo-profiles/get-all-cpu.pprof &
```

Запустить запросы:

```bash
seq 1 200 | xargs -P 20 -I '{}' \
  curl -fsS -o /dev/null \
  -H 'X-User-ID: 1' \
  http://localhost:8080/api/v1/todos
```

После завершения подождать 30–60 секунд и снять heap:

```bash
curl -fsS \
  'http://localhost:9090/debug/pprof/heap?gc=1' \
  -o /tmp/go-todo-profiles/heap-after-get-all.pprof
```

Сравнить:

```bash
go tool pprof \
  -base /tmp/go-todo-profiles/heap-baseline.pprof \
  /tmp/go-todo-profiles/heap-after-get-all.pprof
```

В интерактивном режиме:

```text
top
top -cum
list todoRepository.GetAll
list TodoHandler.GetAll
```

Ожидаемая проблема: большой кратковременный рост памяти из-за загрузки всего результата из БД и последующей JSON-сериализации. Если память после GC возвращается к baseline, это memory spike, но не утечка.

## 8. Проверить goroutine leak

После нагрузки:

```bash
curl -fsS \
  'http://localhost:9090/debug/pprof/goroutine?debug=2' \
  -o /tmp/go-todo-profiles/goroutine-after.txt
```

Сравнить количество goroutine:

```bash
rg -c '^goroutine [0-9]+' \
  /tmp/go-todo-profiles/goroutine-baseline.txt \
  /tmp/go-todo-profiles/goroutine-after.txt
```

Посмотреть повторяющиеся состояния:

```bash
rg '\[chan send|\[chan receive|\[select|\[IO wait' \
  /tmp/go-todo-profiles/goroutine-after.txt
```

Признак утечки — одинаковые пользовательские стеки, число которых растёт после каждого цикла и не уменьшается в покое.

HTTP connection goroutine в `IO wait` непосредственно во время нагрузки сами по себе утечкой не являются.

## 9. Проверить mutex и block contention

В compose sampling уже включён.

Под нагрузкой снять профили:

```bash
curl -fsS \
  http://localhost:9090/debug/pprof/mutex \
  -o /tmp/go-todo-profiles/mutex.pprof

curl -fsS \
  http://localhost:9090/debug/pprof/block \
  -o /tmp/go-todo-profiles/block.pprof
```

Анализ:

```bash
go tool pprof -top /tmp/go-todo-profiles/mutex.pprof
go tool pprof -top /tmp/go-todo-profiles/block.pprof
```

Искать:

- contention в Zap/Lumberjack;
- блокировки пула БД;
- gRPC stream waits;
- channel send/receive;
- runtime mutex contention.

## 10. Проверить execution trace

Во время нагрузки:

```bash
curl -fsS \
  'http://localhost:9090/debug/pprof/trace?seconds=10' \
  -o /tmp/go-todo-profiles/trace.out
```

Открыть:

```bash
go tool trace /tmp/go-todo-profiles/trace.out
```

Проверить:

- goroutine analysis;
- network blocking;
- synchronization blocking;
- GC pauses;
- scheduler latency;
- чрезмерное создание goroutine.

## 11. Критерии проблемы

Утечка памяти подтверждается, если после нескольких циклов нагрузка → покой → GC:

- `heap_inuse` растёт ступенчато;
- `heap_objects` не возвращается;
- heap diff показывает один и тот же удерживающий стек;
- рост не объясняется увеличением постоянного кэша.

Goroutine leak подтверждается, если:

- `go_goroutines` растёт после каждого цикла;
- dump показывает одинаковые пользовательские стеки;
- goroutine остаются после отмены запросов.

Проблема пула БД подтверждается, если:

- in-use соединения достигают максимума;
- растут `wait_count` и `wait_duration`;
- SQL spans становятся медленными;
- одновременно растут HTTP/gRPC latency.

## 12. Завершить эксперимент

Остановить сервисы, сохранив БД:

```bash
docker compose down
```

Удалить также тестовые данные PostgreSQL:

```bash
docker compose down -v
```

Первым практическим кандидатом рекомендую исследовать `GetAll` с 10–100 тысячами записей: он проще всего воспроизводится через REST и хорошо виден одновременно в Prometheus, логах, Jaeger, heap, CPU и execution trace.