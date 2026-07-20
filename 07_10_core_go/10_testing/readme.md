WordCount - Подсчитывает частоту слов, игнорируя регистр и разделители
RenderReport - Формирует отсортированный текстовый отчёт из карты подсчётов
NormalizeSpaces - Убирает лишние пробелы и нормализует строку

все
```bash
go test ./...
```
```bash
go test .
```

из пакета
```bash
go test ./testingdemo
```

по имени
```bash
go test -run TestWordCount
```

```bash
go test -run WordCount
```

запуск фаззинг теста
go test - стандартный запуск тестов
-fuzz=FuzzWordCount - тот флаг переключает go test из режима обычных юнит-тестов в режим фаззера
-fuzztime=10s - ограничивает время работы фаззера (Вместо секунд можно использовать x, например -fuzztime=1000x, чтобы задать точное количество итераций, но секунды используются чаще)
```bash
go test -fuzz=FuzzWordCount -fuzztime=10s
```