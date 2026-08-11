# Примеры к теме 23: память и GC

```bash
go run ./cmd/stackheap
go build -gcflags='go_examples/23_memory/memdemo=-m=2' ./cmd/stackheap
go test -bench=. -benchmem ./memdemo

go run ./cmd/gcstats -blocks=64 -mb=1
GODEBUG=gctrace=1 go run ./cmd/gcstats -blocks=256 -mb=1
GOGC=50 GODEBUG=gctrace=1 go run ./cmd/gcstats -blocks=256 -mb=1
```

`runtime.GC` в `gcstats` нужен только для воспроизводимого учебного опыта. Прикладная программа обычно позволяет runtime самостоятельно выбирать момент сборки.
