## Строки, байты и руны (25 мин)

### 3.1. Строка — неизменяемая последовательность байт
Строка в Go всегда в кодировке **UTF‑8**.  
Длина в байтах: `len(s)`  
Длина в символах (рунах): `utf8.RuneCountInString(s)` или `range s`.

### 3.2. Рунные литералы и итерация
```go
s := "Привет"
fmt.Println(len(s))                // 12 (байт)
fmt.Println(utf8.RuneCountInString(s)) // 6
for i, r := range s {              // i — смещение в байтах, r — руна
    fmt.Printf("%d: %c %U\n", i, r, r)
}
```

### 3.3. Доступ к первому байту и первой руне
```go
func firstRune(s string) (rune, int) {
    return utf8.DecodeRuneInString(s)
}
```
Если строка пуста, `DecodeRuneInString` возвращает `(0, 0)`.

### 3.4. Практический пример: AnalyseString
```go
func AnalyzeString(s string) (byteLen, runeCount, firstByte, firstRune int, err error) {
    if s == "" {
        return 0, 0, 0, 0, errors.New("empty string")
    }
    byteLen = len(s)
    runeCount = utf8.RuneCountInString(s)
    firstByte = int(s[0])
    firstRune, _ = utf8.DecodeRuneInString(s)
    return
}
```
