# main.go

Короткое введение в пакет `reflect`: тип, значение, изменение и обход структуры.

Показывает:

- `reflect.TypeOf` и `Kind` для `int` и `string`;
- `reflect.ValueOf`, `Interface`, `Int`, `String`;
- изменение переменной через указатель: `ValueOf(&x).Elem().SetInt(...)`;
- обход полей структуры `Person` через `NumField` / `Field(i)`.

Команда:

```bash
go run main.go
```

Не запускайте `go run .`: в каталоге несколько независимых `main`.
