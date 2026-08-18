# 01_example.go

Минимальный пример рефлексии и ловушки typed nil.

Показывает:

- `reflect.TypeOf` / `reflect.ValueOf` для целого числа;
- различие `nil`-интерфейса и интерфейса с typed nil (`*int`, `*MyError`);
- `Type.Name`, `PkgPath`, `Kind` для именованного типа `UserID`;
- простую функцию `describe` с ветвлением по `Kind`.

Команда:

```bash
go run 01_example.go
```

Не запускайте `go run .`: в каталоге несколько независимых `main`.
