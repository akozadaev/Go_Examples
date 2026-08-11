# RabbitMQ example for Go

Учебный проект к лекции об очередях сообщений. Он показывает полную цепочку
RabbitMQ: topic exchange, quorum queues, persistent messages, publisher
confirms, manual acknowledgements, prefetch, delayed retry, DLQ, подавление
дубликатов в пределах процесса и graceful shutdown.

Для main и retry quorum queues включена стратегия at-least-once dead-lettering:
исходная очередь удерживает dead-lettered сообщение до publisher confirm от
целевой очереди. Требуемый этой стратегией overflow — `reject-publish`.

## Архитектура

```text
publisher --orders.created--> orders.events (topic)
                                  |
                                  v
                         orders.created.worker
                                  |
                                  v
                               consumer

transient error -> confirmed republish -> orders.retry
                                            |
                                      retry queue (TTL)
                                            |
                                            +--> orders.events

permanent error / retries exhausted -> orders.dlx -> orders.created.dead
```

При временной ошибке consumer сначала подтверждённо публикует копию в retry
exchange и только после publisher confirm делает ack исходного delivery. Если
retry publish не подтверждён, оригинал остаётся unacked и возвращается RabbitMQ
после закрытия channel. Если publish подтверждён, а ack потерян, возможен
дубликат — это нормальная семантика at-least-once.

Publisher использует `mandatory=true`: confirm доказывает приём публикации
RabbitMQ, а `basic.return` отдельно выявляет отсутствие подходящего binding.

## Требования

- Go 1.26.4;
- Docker с Compose;
- RabbitMQ 4.2 Management из `compose.yaml`.

## Быстрый запуск

```bash
docker compose up -d --wait
go run ./cmd/consumer
```

В другом терминале опубликуйте событие:

```bash
go run ./cmd/publisher \
  --order-id order-42 \
  --customer-id customer-7 \
  --amount 159900 \
  --currency RUB
```

Management UI: <http://localhost:15672>. Логин и пароль локального окружения:
`app` / `app`.

Приложения читают конфигурацию непосредственно из environment. Значения и
локальные defaults перечислены в `.env.example`; файл `.env` автоматически не
загружается.

## Демонстрационные сценарии

Успешная обработка:

```bash
go run ./cmd/publisher --order-id order-success
```

Две временные ошибки, затем успешная обработка после двух задержек:

```bash
go run ./cmd/publisher \
  --order-id order-retry \
  --demo-failure transient \
  --fail-attempts 2
```

Постоянная ошибка сразу отправляет delivery в DLQ:

```bash
go run ./cmd/publisher \
  --order-id order-permanent \
  --demo-failure permanent
```

Повреждённый JSON также попадает в DLQ:

```bash
go run ./cmd/publisher --order-id order-malformed --malformed
```

Если `--fail-attempts` больше `MAX_RETRIES`, сообщение попадает в DLQ после
исчерпания разрешённых повторов:

```bash
MAX_RETRIES=2 go run ./cmd/publisher \
  --order-id order-exhausted \
  --demo-failure transient \
  --fail-attempts 10
```

`MAX_RETRIES` является настройкой consumer, поэтому consumer тоже должен быть
запущен с нужным значением.

Технические headers `x-demo-failure` и `x-demo-fail-attempts` существуют только
для воспроизводимой демонстрации. Они не входят в бизнес-контракт JSON.

## Просмотр DLQ

Без `--ack` команда показывает одно сообщение и возвращает его в DLQ:

```bash
go run ./cmd/dlq-reader
```

У DLQ отключён стандартный quorum delivery limit, чтобы повторный безопасный
просмотр не удалил сообщение после 20 requeue. Это осознанное исключение только
для очереди ручного разбора; main queue сохраняет защиту RabbitMQ от poison
requeue loop.

Чтобы после просмотра удалить сообщение:

```bash
go run ./cmd/dlq-reader --ack
```

В production просмотр и replay DLQ требуют авторизации, аудита, redaction
чувствительных данных и отдельного runbook. Этот CLI намеренно не реализует
автоматический replay.

## Дубликаты и идемпотентность

Флаг `--message-id` позволяет повторно отправить один логический event ID:

```bash
go run ./cmd/publisher --order-id order-duplicate --message-id demo-event-1
go run ./cmd/publisher --order-id order-duplicate --message-id demo-event-1
```

Consumer подавляет повтор в пределах жизни процесса с помощью
`MemoryProcessedStore`. Это только учебная демонстрация. После перезапуска
consumer память очищается.

Production-обработчик должен в одной транзакции:

1. вставить `(consumer_name, message_id)` в таблицу `processed_messages`;
2. выполнить бизнес-изменение;
3. сделать commit;
4. только после commit отправить RabbitMQ ack.

## Конфигурация

| Переменная | Default | Назначение |
|---|---:|---|
| `AMQP_URL` | `amqp://app:app@localhost:5672/app` | Адрес RabbitMQ |
| `PREFETCH` | `8` | Максимум unacked deliveries |
| `RETRY_DELAY` | `5s` | TTL retry queue |
| `MAX_RETRIES` | `3` | Число повторов после первой попытки |
| `PUBLISH_TIMEOUT` | `10s` | Timeout публикации и confirm |
| `SHUTDOWN_TIMEOUT` | `10s` | Максимальное ожидание остановки consumer |

`RETRY_DELAY` является immutable argument уже созданной очереди. После изменения
значения существующую `orders.created.retry` нужно мигрировать или удалить в
локальном окружении. Полное удаление локальных данных:

```bash
docker compose down -v
```

Команда удаляет volume RabbitMQ вместе со всеми локальными сообщениями.

## Проверки

```bash
go test ./...
go test -race ./...
go build ./...

# Требует запущенный RabbitMQ
go test -tags=integration ./internal/messaging
```

Unit-тесты не требуют RabbitMQ. Для интеграционной проверки используется
настоящий broker из `compose.yaml`. Интеграционный тест также проверяет, что
mandatory publication без binding возвращается publisher как ошибка.

## Важные ограничения

- Соединения автоматически не восстанавливаются. При сетевой ошибке процесс
  завершается с ошибкой и должен быть перезапущен supervisor/orchestrator.
- Confirm timeout имеет неопределённый результат: RabbitMQ мог принять сообщение,
  даже если приложение не получило confirm.
- Один `Publisher` сериализует публикации и ждёт каждый confirm. Это корректно и
  понятно для обучения, но не предназначено для высокого throughput.
- DLX и TTL заданы queue arguments. В production RabbitMQ рекомендует policies,
  когда настройку можно вынести из приложения.
- Quorum queue на одном локальном узле не обеспечивает отказоустойчивость;
  production-кластеру требуется корректное нечётное число replicas.
- Учебные credentials и management port нельзя переносить в production.
