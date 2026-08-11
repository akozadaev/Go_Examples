# 27. Профилирование и оптимизация Go-программ

**Продолжительность:** 2 академических часа  
**Примеры:** `27_profiling` и `21_tracing/go_todo_service`

## Цели

После лекции студент сможет:

- формулировать измеримую гипотезу оптимизации;
- писать воспроизводимые Go-бенчмарки;
- интерпретировать `ns/op`, `B/op` и `allocs/op`;
- снимать CPU, heap, allocs, goroutine, mutex и block-профили;
- читать `top`, `list`, web-граф и flame graph;
- отличать tracing, metrics и profiling;
- доказательно сравнивать состояние до и после изменения;
- безопасно эксплуатировать `pprof` в сервисе.

## 1. Оптимизация как эксперимент

Правильный цикл:

```text
цель → измерение baseline → гипотеза → одно изменение
  ↑                                      │
  └──── проверка корректности ← повторное измерение
```

Сначала определяют пользовательскую цель: latency, throughput, CPU, память, стоимость или стабильность. Быстрый микробенчмарк бесполезен, если оптимизируется незначимый участок.

Порядок действий:

1. Зафиксировать workload и среду.
2. Проверить корректность и тесты.
3. Получить baseline.
4. Найти hotspot профилем.
5. Сформулировать причинную гипотезу.
6. Сделать минимальное изменение.
7. Сравнить результаты статистически.
8. Проверить системные метрики и регрессии.

## 2. Бенчмарки Go

Benchmark располагается в `_test.go` и получает `*testing.B`:

```go
func BenchmarkFast20(b *testing.B) {
    for b.Loop() {
        result = fib.Fast(20)
    }
}
```

В актуальном Go `b.Loop()` управляет числом итераций и корректно взаимодействует с таймером. В старом стиле используется `for i := 0; i < b.N; i++`.

Мини-пример сравнивает намеренно медленный рекурсивный и линейный алгоритмы:

```bash
cd 27_profiling
go test ./...
go test -bench=. -benchmem -count=5 ./fib
```

Результат вида:

```text
BenchmarkFast20-8    100000000    12.3 ns/op    0 B/op    0 allocs/op
```

означает имя, `GOMAXPROCS`, число операций, время операции, выделенные байты и heap-аллокации.

### 2.1 Защита от оптимизации

Если результат не наблюдаем, компилятор может удалить работу. Поэтому пример сохраняет его в package-level переменную. Но такая переменная тоже может влиять на escape analysis. Бенчмарк должен моделировать реальное использование результата.

### 2.2 Setup и таймер

Подготовку не включают в измерение:

```go
func BenchmarkParse(b *testing.B) {
    input := loadFixture()
    b.ReportAllocs()
    b.ResetTimer()
    for b.Loop() {
        parse(input)
    }
}
```

Для части цикла доступны `b.StopTimer()` и `b.StartTimer()`, но частое переключение само добавляет шум.

### 2.3 Sub-benchmarks

```go
for _, size := range []int{10, 100, 1000} {
    b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
        // Бенчмарк для текущего размера.
    })
}
```

Размер входа — часть контракта теста. Алгоритмы могут менять относительную эффективность при росте данных.

### 2.4 Параллельные бенчмарки

```go
b.RunParallel(func(pb *testing.PB) {
    for pb.Next() {
        sharedOperation()
    }
})
```

Они полезны для конкурентного cache, pool или handler, но измеряют другой вопрос, чем последовательный benchmark. Проверяйте также `go test -race`.

## 3. Надёжность измерения

Шум создают:

- энергосбережение и turbo boost CPU;
- другие процессы;
- фоновые обновления и контейнеры;
- разогрев cache;
- сеть и удалённая база;
- разные версии Go и зависимости;
- изменение `GOMAXPROCS`.

Поэтому фиксируют окружение и выполняют несколько запусков:

```bash
go test -run='^$' -bench=. -benchmem -count=10 ./fib > before.txt
# изменение
go test -run='^$' -bench=. -benchmem -count=10 ./fib > after.txt
benchstat before.txt after.txt
```

Смотрите не только проценты, но и доверие к различию, абсолютный эффект и соответствие production workload.

Нельзя писать тест, который утверждает `Fast` быстрее `Slow` по wall clock. Такие проверки нестабильны. Корректность проверяет unit test, производительность — benchmark и статистическое сравнение.

## 4. Профилирование

Benchmark показывает **сколько**, profile помогает понять **где**.

| Профиль | Главный вопрос |
|---|---|
| CPU | где процесс проводит CPU-время? |
| heap | какие живые объекты удерживают память? |
| allocs | где выделялась память за всё время? |
| goroutine | какие горутины существуют и где ждут? |
| mutex | где ожидали захваченный mutex? |
| block | где блокировались на синхронизации? |
| threadcreate | почему создавались потоки ОС? |
| execution trace | как взаимодействуют scheduler, goroutines, syscalls и GC? |

Профиль является выборкой, а не полным журналом каждого события. Короткий или ненагруженный CPU-профиль может быть пустым.

## 5. Профили из benchmark

```bash
go test -run='^$' -bench=Slow \
  -cpuprofile=cpu.pprof \
  -memprofile=mem.pprof \
  ./fib
```

Полезно увеличить длительность:

```bash
go test -run='^$' -bench=Slow -benchtime=10s -cpuprofile=cpu.pprof ./fib
```

Просмотр в стандартной поставке Go:

```bash
go tool pprof -top cpu.pprof
go tool pprof -list='fib\.Slow' cpu.pprof
go tool pprof -http=:8081 cpu.pprof
```

В интерактивном режиме:

- `top` — верхние узлы;
- `top -cum` — сортировка по cumulative;
- `list Function` — строки исходника;
- `peek regexp` — callers/callees;
- `web` — граф вызовов, обычно нужен Graphviz.

`flat` — стоимость непосредственно в функции. `cum` — функция вместе со всем вызванным ею кодом. Высокий `cum` при низком `flat` означает, что работа выполняется глубже.

## 6. HTTP pprof

Пакет `net/http/pprof` регистрирует обработчики на `http.DefaultServeMux`:

```go
import _ "net/http/pprof"

log.Fatal(http.ListenAndServe("127.0.0.1:6060", nil))
```

Мини-сервис находится в [`cmd/pprofdemo`](./cmd/pprofdemo/main.go):

```bash
go run ./cmd/pprofdemo
curl 'http://127.0.0.1:6060/work?n=35'
```

CPU-профиль снимается за интервал, в течение которого нужна нагрузка:

```bash
curl -o cpu.pprof \
  'http://127.0.0.1:6060/debug/pprof/profile?seconds=15'
```

Heap — моментальный профиль:

```bash
curl -o heap.pprof 'http://127.0.0.1:6060/debug/pprof/heap'
curl -o allocs.pprof 'http://127.0.0.1:6060/debug/pprof/allocs'
```

### Безопасность

Pprof раскрывает стеки, имена функций, аргументы запуска и детали поведения системы. CPU-profile и trace создают дополнительную нагрузку. Поэтому endpoint:

- не публикуют в интернет;
- привязывают к loopback или отдельному admin-интерфейсу;
- защищают network policy, VPN или аутентификацией;
- ограничивают по времени;
- не смешивают с публичным API без осознанной защиты.

## 7. Heap и allocs

Heap-profile имеет разные sample types:

- `inuse_space` — байты живых объектов;
- `inuse_objects` — число живых объектов;
- `alloc_space` — суммарно выделенные байты;
- `alloc_objects` — суммарное число аллокаций.

Примеры:

```bash
go tool pprof -sample_index=inuse_space heap.pprof
go tool pprof -sample_index=alloc_space allocs.pprof
```

Если нужно снизить давление на GC, часто начинают с `alloc_space` и `alloc_objects`. Если процесс удерживает всё больше памяти, сравнивают `inuse_space` в разные моменты.

Профили можно сравнивать:

```bash
go tool pprof -base=before.pprof after.pprof
```

Но workload обоих снимков должен быть сопоставим.

## 8. Mutex и block profiles

Эти профили не всегда собирают с полной частотой по умолчанию. Для учебного сервиса sampling включают осознанно:

```go
runtime.SetMutexProfileFraction(1)
runtime.SetBlockProfileRate(1)
```

Значение `1` удобно для короткой демонстрации, но может быть слишком дорогим для production. После эксперимента настройки возвращают или процесс перезапускают.

Высокое mutex contention не обязательно означает, что надо заменить mutex. Возможны слишком большая критическая секция, один глобальный lock, медленный вызов под lock или чрезмерная конкурентность.

## 9. Execution trace

Trace полезен, когда CPU-profile недостаточен:

- задержки scheduler;
- массовое создание горутин;
- блокировки и syscalls;
- влияние GC;
- распределение работы между процессорами.

Для теста:

```bash
go test -trace=trace.out ./...
go tool trace trace.out
```

Для HTTP pprof:

```bash
curl -o trace.out \
  'http://127.0.0.1:6060/debug/pprof/trace?seconds=5'
```

Trace быстро растёт и сам влияет на процесс, поэтому интервал должен быть коротким и содержательным.

## 10. Todo-сервис как стенд

В [`21_tracing/go_todo_service`](../21_tracing/go_todo_service) уже подключены pprof-маршруты и OpenTelemetry. Добавлены:

- benchmarks модели и user context;
- `scripts/test_profiling.sh` для базовых профилей;
- `scripts/run_profiling_test.sh` для нагрузки и повторного снимка;
- Makefile-цели `profile` и `profile-load`.

Benchmark без инфраструктуры:

```bash
cd 21_tracing/go_todo_service
go test -bench=. -benchmem ./internal/model ./internal/userctx
```

После запуска сервиса в `GIN_MODE=debug`:

```bash
make profile
PROFILE_REQUESTS=5000 PROFILE_CONCURRENCY=50 make profile-load
```

CPU-профилирование должно идти одновременно с нагрузкой. Скрипт сначала прогревает endpoint, затем запускает CPU-profile в фоне и повторяет нагрузку.

`pprof` и OpenTelemetry решают разные задачи:

- trace показывает путь и latency конкретного запроса между компонентами;
- profile агрегирует использование ресурсов всего процесса;
- metrics показывают изменение чисел во времени.

## 11. Поиск оптимизации

Пример процесса:

1. SLO: p95 latency `GET /todos` меньше 100 мс.
2. Нагрузка показывает 180 мс.
3. Trace показывает большую часть времени в БД — CPU-оптимизация JSON не поможет.
4. Если CPU-profile показывает сериализацию, исследуем её `flat/cum`.
5. Allocs-profile показывает формирование временных слайсов логгера.
6. Пишем изолированный benchmark этой операции.
7. Меняем только реализацию формирования полей.
8. Сравниваем `benchstat`, CPU/alloc profiles и end-to-end latency.

Профиль указывает место затрат, но не готовое решение. Причину нужно подтвердить кодом и экспериментом.

## 12. Типичные ошибки

- Профилировать сервис без репрезентативной нагрузки.
- Оптимизировать функцию только из-за высокого `cum`, не посмотрев callees.
- Сравнивать профили разных workload.
- Делать вывод по одному benchmark-запуску.
- Включать файловый I/O, сеть или setup в микробенчмарк случайно.
- Считать снижение `ns/op` достаточным, игнорируя память и корректность.
- Публиковать pprof endpoint наружу.
- Снимать слишком короткий CPU-профиль.
- Оптимизировать сгенерированный или библиотечный код до анализа собственного вызова.
- Менять несколько факторов одновременно.

## 13. Контрольные вопросы

1. Чем benchmark отличается от load test?
2. Зачем запускать benchmark несколько раз?
3. Что означают `B/op` и `allocs/op`?
4. Почему результат benchmark нужно сделать наблюдаемым?
5. Чем `flat` отличается от `cum`?
6. Чем `heap` отличается от `allocs`?
7. Почему CPU-profile без нагрузки пуст?
8. Когда нужен execution trace?
9. Почему pprof endpoint нельзя публиковать?
10. Как доказать, что оптимизация улучшила продуктовую метрику?

## Итоги

- Оптимизируют измеримую проблему, а не внешний вид кода.
- Benchmark нужен для контролируемого сравнения небольшой операции.
- Pprof локализует CPU, memory и contention hotspots.
- Профили снимают под репрезентативной и повторяемой нагрузкой.
- Результат подтверждают повторным измерением, тестами и системной метрикой.
- Tracing, profiling и metrics дополняют друг друга.

## Источники

- [Go diagnostics](https://go.dev/doc/diagnostics)
- [`testing.B`](https://pkg.go.dev/testing#B)
- [`runtime/pprof`](https://pkg.go.dev/runtime/pprof)
- [`net/http/pprof`](https://pkg.go.dev/net/http/pprof)
- [Profiling Go Programs](https://go.dev/blog/pprof)
- [Execution tracer](https://go.dev/blog/execution-traces-2024)
