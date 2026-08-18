# main_valueof.go

Работа с `reflect.Value`: чтение, указатели и изменение значения.

Показывает:

- `ValueOf`, `Interface`, `Type`, `Kind`;
- почему `ValueOf(x)` не settable, а `ValueOf(&x).Elem()` - да;
- `Elem` для разыменования указателя;
- изменение исходной переменной через `SetInt`.

Команда:

```bash
go run main_valueof.go
```

Не запускайте `go run .`: в каталоге несколько независимых `main`.
