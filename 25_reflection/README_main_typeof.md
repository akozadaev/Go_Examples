# main_typeof.go

Демонстрация `reflect.Type` и `TypeOf` для разных категорий типов.

Показывает:

- получение точного типа и базовой категории `Kind`;
- примеры для `int`, `string`, `bool`, `float64`;
- `Kind` для слайса, map и структуры;
- отличие имени типа (`Person`) от kind (`struct`) и kind слайса (`slice`).

Команда:

```bash
go run main_typeof.go
```

Не запускайте `go run .`: в каталоге несколько независимых `main`.
