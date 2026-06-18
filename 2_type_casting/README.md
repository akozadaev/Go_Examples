## Безопасное приведение целых чисел (20 мин)

### 2.1. Проверка диапазона
Используйте константы из пакета `math`:
- `math.MaxInt8`, `math.MinInt8`, `math.MaxInt32`, `math.MinInt32`, `math.MaxInt64`, `math.MinInt64`, аналогично для беззнаковых.

**Шаблон безопасного приведения:**
```go
func SafeInt32ToInt8(v int32) (int8, error) {
    if v > math.MaxInt8 || v < math.MinInt8 {
        return 0, fmt.Errorf("value %d out of int8 range", v)
    }
    return int8(v), nil
}
```

### 2.2. Приведение float → int
Приведение `float64` к целому **усекает** дробную часть (не округляет).  
Дополнительно нужно проверять `NaN` и `Inf`.

```go
func SafeFloatToInt64(f float64) (int64, error) {
    if math.IsNaN(f) || math.IsInf(f, 0) {
        return 0, fmt.Errorf("invalid float: %v", f)
    }
    if f > math.MaxInt64 || f < math.MinInt64 {
        return 0, fmt.Errorf("out of int64 range")
    }
    return int64(f), nil
}
```

### 2.3. Примеры с табличными тестами
```go
func TestSafeInt32ToInt8(t *testing.T) {
    tests := []struct{
        name string
        v int32
        want int8
        wantErr bool
    }{
        {"max", 127, 127, false},
        {"min", -128, -128, false},
        {"overflow", 128, 0, true},
        {"underflow", -129, 0, true},
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := SafeInt32ToInt8(tt.v)
            if (err != nil) != tt.wantErr { ... }
            if got != tt.want { ... }
        })
    }
}
```