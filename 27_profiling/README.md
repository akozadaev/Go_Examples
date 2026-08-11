# Примеры к теме 27: benchmarks и pprof

```bash
go test ./...
go test -bench=. -benchmem -count=5 ./fib
go test -bench=. -benchmem -cpuprofile=cpu.pprof -memprofile=mem.pprof ./fib
```

Мини-сервер `pprof` привязан только к loopback-интерфейсу:

```bash
go run ./cmd/pprofdemo
curl 'http://127.0.0.1:6060/work?n=35'
curl -o cpu.pprof 'http://127.0.0.1:6060/debug/pprof/profile?seconds=15'
curl -o heap.pprof 'http://127.0.0.1:6060/debug/pprof/heap'
```

Пока снимается CPU-профиль, создайте нагрузку в другом терминале:

```bash
for i in $(seq 1 100); do curl -s 'http://127.0.0.1:6060/work?n=35' >/dev/null; done
```

Просмотр профиля зависит от комплектации Go toolchain. В стандартной поставке:

```bash
go tool pprof -top cpu.pprof
go tool pprof -http=:8081 cpu.pprof
```
