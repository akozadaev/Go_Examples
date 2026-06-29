# Лекция: Управляющие конструкции языка и отладка программ в Go

**Длительность:** 2 астрономических часа (120 минут)  
**Целевая аудитория:** разработчики, знакомые с базовым синтаксисом Go: переменные, функции, пакеты, простые типы.  
**Формат:** объяснение, живое кодирование, короткие вопросы аудитории, демонстрация отладки.

---

## Цели лекции

После занятия слушатель должен уметь:

- применять `if`, `else`, `switch` и `type switch`;
- понимать, в каком порядке выполняются условия, ветки, циклы и выражения;
- объяснять короткие замыкания `&&` и `||`;
- безопасно проверять `nil` с помощью коротких замыканий;
- понимать, когда вычисляются аргументы функций;
- писать все основные формы цикла `for`;
- использовать `break`, `continue` и метки;
- выполнять базовую отладку через `go run`, `go build` и Delve (`dlv`).

---

## Хронометраж на 120 минут

| Блок | Время |
|------|-------|
| Введение, цель темы, связь с предыдущими лекциями | 5 мин |
| `if/else`: синтаксис, короткая инструкция, область видимости | 20 мин |
| Порядок выполнения и вычисления выражений | 15 мин |
| Логические операции, короткие замыкания, `nil`-проверки | 20 мин |
| `switch`, `fallthrough`, `type switch` | 20 мин |
| Цикл `for`: формы, `range`, порядок выполнения | 20 мин |
| Метки, `break`, `continue` во вложенных циклах | 8 мин |
| Отладка управляющих конструкций | 10 мин |
| Итоги и вопросы | 2 мин |
| **Всего** | **120 мин** |

---

## Материалы для преподавателя

- Установленный Go.
- Установленный Delve: `go install github.com/go-delve/delve/cmd/dlv@latest`.
- Репозиторий с примерами занятия.
- Официальные страницы для сверки:
  - [Go Tour: Flow control](https://go.dev/tour/flowcontrol/1)
  - [Go Spec: If statements](https://go.dev/ref/spec#If_statements)
  - [Go Spec: Switch statements](https://go.dev/ref/spec#Switch_statements)
  - [Go Spec: For statements](https://go.dev/ref/spec#For_statements)
  - [Go Spec: Order of evaluation](https://go.dev/ref/spec#Order_of_evaluation)

---

# 1. `if/else`: условия и ветвление

## 1.1. Базовая форма `if`

В Go условие не заключается в круглые скобки, но тело ветки всегда пишется в фигурных скобках.

```go
package main

import "fmt"

func main() {
    x := -5

    if x > 0 {
        fmt.Println("positive")
    } else {
        fmt.Println("non-positive")
    }
}
```

Ключевые правила:

- условие должно иметь тип `bool`;
- `if x {}` не скомпилируется, если `x` имеет тип `int`;
- `else` пишется на той же строке, где закрывается блок `if`;
- лишние круглые скобки вокруг условия не нужны.

Пример ошибки:

```go
var x int = 1

if x {
    fmt.Println("will not compile")
}
```

Go не приводит числа к `bool` автоматически. Это осознанное отличие от C-подобных языков.

---

## 1.2. `if` с короткой инструкцией

Перед условием можно выполнить короткую инструкцию. Чаще всего это используется для получения значения и проверки ошибки.

```go
package main

import (
    "fmt"
    "os"
)

func main() {
    if file, err := os.Open("main.go"); err != nil {
        fmt.Println("open error:", err)
    } else {
        defer file.Close()
        fmt.Println("file opened")
    }
}
```

Важно:

- сначала выполняется короткая инструкция `file, err := os.Open("main.go")`;
- затем проверяется условие `err != nil`;
- переменные `file` и `err` видны только внутри `if` и `else`;
- после завершения конструкции `if/else` эти переменные недоступны.

Такой стиль делает код компактнее и уменьшает область видимости временных переменных.

---

## 1.3. Ранний выход и "счастливый путь"

В Go часто используют ранний `return`, чтобы не увеличивать вложенность.

```go
package main

import "fmt"

func enter(age int, hasPermission bool) {
    if age < 18 && !hasPermission {
        fmt.Println("access denied")
        return
    }

    fmt.Println("welcome")
}

func main() {
    enter(16, false)
    enter(20, false)
}
```

Идея: сначала отсекаем ошибочные или запрещенные сценарии, затем оставляем основной сценарий на верхнем уровне функции.

---

# 2. Порядок выполнения и вычисления выражений

Этот блок важен для понимания отладки: программа выполняет не "примерно то, что написано", а строго определенную последовательность действий.

## 2.1. Порядок выполнения `if`

Для конструкции:

```go
if init; condition {
    body
} else {
    elseBody
}
```

порядок такой:

1. Выполняется `init`, если он есть.
2. Вычисляется `condition`.
3. Если условие равно `true`, выполняется блок `body`.
4. Если условие равно `false`, выполняется блок `elseBody`, если он есть.
5. После выбранной ветки выполнение продолжается после всей конструкции `if/else`.

Демонстрация:

```go
package main

import "fmt"

func main() {
    if x := 10; x > 5 {
        fmt.Println("branch if")
    } else {
        fmt.Println("branch else")
    }

    fmt.Println("after if")
}
```

---

## 2.2. Порядок вычисления в выражениях

В Go важно различать:

- порядок выполнения инструкций;
- порядок вычисления частей выражения;
- короткие замыкания логических операторов.

Функциональные вызовы, вызовы методов и операции получения из канала внутри выражения вычисляются в лексическом порядке: слева направо.

```go
package main

import "fmt"

func value(name string, result int) int {
    fmt.Println("call", name)
    return result
}

func main() {
    x := value("A", 1) + value("B", 2) + value("C", 3)
    fmt.Println("x =", x)
}
```

Вывод:

```text
call A
call B
call C
x = 6
```

Практический вывод: если выражение содержит вызовы функций с побочными эффектами, порядок этих вызовов важен. Но лучше не писать код, где понимание результата зависит от сложной смеси побочных эффектов внутри одного выражения.

---

## 2.3. Вычисление аргументов функций

Перед входом в функцию Go сначала вычисляет значение функции и все аргументы вызова.

```go
package main

import "fmt"

func arg(name string, value int) int {
    fmt.Println("arg", name)
    return value
}

func sum(a, b int) int {
    fmt.Println("inside sum")
    return a + b
}

func main() {
    result := sum(arg("first", 10), arg("second", 20))
    fmt.Println("result =", result)
}
```

Вывод:

```text
arg first
arg second
inside sum
result = 30
```

Что важно проговорить:

- аргументы вычисляются до начала выполнения тела функции;
- если вычисление аргумента вызывает панику, функция не начнет выполняться;
- если функция не использует один из параметров, переданный аргумент все равно уже был вычислен;
- для `defer` аргументы тоже вычисляются сразу, а сам вызов откладывается.

Пример с `defer`:

```go
package main

import "fmt"

func main() {
    x := 1
    defer fmt.Println("defer:", x)

    x = 2
    fmt.Println("main:", x)
}
```

Вывод:

```text
main: 2
defer: 1
```

`fmt.Println("defer:", x)` будет выполнен в конце функции, но значение `x` для аргумента было вычислено в момент регистрации `defer`.

---

# 3. Логические операции и короткие замыкания

## 3.1. Основные логические операторы

В Go есть три основных логических оператора:

| Оператор | Значение | Пример |
|----------|----------|--------|
| `&&` | логическое И | `age >= 18 && hasID` |
| `||` | логическое ИЛИ | `isAdmin || isOwner` |
| `!` | логическое НЕ | `!blocked` |

Приоритет:

1. `!`
2. `&&`
3. `||`

Пример:

```go
if !blocked && age >= 18 || isAdmin {
    fmt.Println("allowed")
}
```

Такой код лучше писать с явными скобками:

```go
if (!blocked && age >= 18) || isAdmin {
    fmt.Println("allowed")
}
```

Скобки не меняют работу программы, если приоритет уже такой же, но помогают читать условие.

---

## 3.2. Короткое замыкание `&&`

`&&` вычисляет операнды слева направо. Если левый операнд равен `false`, правый операнд не вычисляется, потому что результат всего выражения уже точно `false`.

```go
package main

import "fmt"

func check() bool {
    fmt.Println("check called")
    return true
}

func main() {
    if false && check() {
        fmt.Println("inside if")
    }

    fmt.Println("after if")
}
```

Вывод:

```text
after if
```

`check()` не вызывается.

Практическое применение:

```go
if user.IsActive && user.HasPaid {
    fmt.Println("show premium content")
}
```

Если `user.IsActive == false`, проверка `user.HasPaid` уже не нужна.

---

## 3.3. Короткое замыкание `||`

`||` тоже вычисляет операнды слева направо. Если левый операнд равен `true`, правый операнд не вычисляется, потому что результат всего выражения уже точно `true`.

```go
package main

import "fmt"

func fallback() bool {
    fmt.Println("fallback called")
    return true
}

func main() {
    if true || fallback() {
        fmt.Println("allowed")
    }
}
```

Вывод:

```text
allowed
```

`fallback()` не вызывается.

Практическое применение:

```go
if isAdmin || isOwner {
    fmt.Println("can edit")
}
```

Если пользователь уже администратор, проверка владельца не нужна.

---

## 3.4. Операторы, у которых нет логического короткого замыкания

Короткое замыкание есть только у логических операторов `&&` и `||`.

Не стоит путать их с побитовыми операторами:

| Оператор | Что делает | Короткое замыкание |
|----------|------------|--------------------|
| `&&` | логическое И | есть |
| `||` | логическое ИЛИ | есть |
| `&` | побитовое И | нет |
| `|` | побитовое ИЛИ | нет |
| `^` | побитовое XOR / NOT | нет |

Пример:

```go
package main

import "fmt"

func left() bool {
    fmt.Println("left")
    return false
}

func right() bool {
    fmt.Println("right")
    return true
}

func main() {
    fmt.Println(left() && right())
}
```

Вывод:

```text
left
false
```

Для `&&` `right()` не вызывается. С побитовыми операторами такая логика не применяется.

---

## 3.5. Безопасная проверка на `nil` через короткое замыкание

Классический пример: перед разыменованием указателя нужно проверить, что он не равен `nil`.

```go
package main

import "fmt"

func main() {
    var p *int

    if p != nil && *p > 0 {
        fmt.Println("positive")
    } else {
        fmt.Println("nil or non-positive")
    }
}
```

Почему это безопасно:

1. Сначала вычисляется `p != nil`.
2. Если `p == nil`, результат левой части `false`.
3. Для `&&` правую часть `*p > 0` уже не нужно вычислять.
4. Разыменования `nil` не происходит, паники нет.

Небезопасный вариант:

```go
if *p > 0 && p != nil {
    fmt.Println("positive")
}
```

Здесь программа сначала попытается вычислить `*p > 0`, и при `p == nil` будет паника.

---

## 3.6. Проверка `nil` у вложенных полей

Короткое замыкание удобно для цепочек проверок.

```go
package main

import "fmt"

type Profile struct {
    Email string
}

type User struct {
    Profile *Profile
}

func main() {
    user := &User{}

    if user != nil && user.Profile != nil && user.Profile.Email != "" {
        fmt.Println("email:", user.Profile.Email)
    } else {
        fmt.Println("email is empty or unavailable")
    }
}
```

Порядок проверок идет от самого внешнего объекта к внутреннему:

1. `user != nil`
2. `user.Profile != nil`
3. `user.Profile.Email != ""`

Если поменять порядок, можно получить панику.

---

## 3.7. Проверка `nil` у map, slice и channel

Не все `nil`-значения одинаково опасны.

```go
package main

import "fmt"

func main() {
    var nums []int
    var dict map[string]int

    fmt.Println(len(nums)) // безопасно: 0
    fmt.Println(len(dict)) // безопасно: 0

    nums = append(nums, 10) // безопасно
    fmt.Println(nums)

    // dict["x"] = 1 // panic: assignment to entry in nil map

    if dict != nil {
        dict["x"] = 1
    }
}
```

Полезно проговорить:

- `nil`-slice можно читать через `len`, и в него можно делать `append`;
- из `nil`-map можно читать, но нельзя записывать;
- отправка в `nil`-channel или чтение из него блокирует горутину;
- короткие замыкания помогают явно защитить опасные операции.

---

# 4. `switch`

## 4.1. Классический `switch`

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    switch day := time.Now().Weekday(); day {
    case time.Saturday, time.Sunday:
        fmt.Println("weekend")
    default:
        fmt.Println("weekday")
    }
}
```

Особенности:

- `switch` может иметь короткую инструкцию;
- выражение после `switch` вычисляется один раз;
- выполняется первый подходящий `case`;
- автоматического проваливания в следующий `case` нет;
- `break` в конце каждого `case` обычно не нужен.

---

## 4.2. Порядок вычисления `switch`

Для выражения:

```go
switch x {
case a:
    // ...
case b, c:
    // ...
default:
    // ...
}
```

порядок такой:

1. Вычисляется `x`.
2. Значения `case` проверяются сверху вниз.
3. Если в одном `case` несколько значений, они проверяются слева направо.
4. Выполняется первый совпавший блок.
5. Остальные `case` не выполняются.

Демонстрация:

```go
package main

import "fmt"

func value(name string, result int) int {
    fmt.Println("value", name)
    return result
}

func main() {
    switch x := value("switch expr", 2); x {
    case value("case 1", 1):
        fmt.Println("one")
    case value("case 2a", 2), value("case 2b", 20):
        fmt.Println("two")
    default:
        fmt.Println("default")
    }
}
```

Ожидаемый вывод:

```text
value switch expr
value case 1
value case 2a
two
```

`case 2b` не вычисляется, потому что совпадение уже найдено на `case 2a`.

---

## 4.3. `switch` без выражения

`switch` можно писать без выражения. Это похоже на `switch true` и удобно для последовательной проверки условий.

```go
package main

import "fmt"

func main() {
    score := 85

    switch {
    case score >= 90:
        fmt.Println("A")
    case score >= 80:
        fmt.Println("B")
    case score >= 70:
        fmt.Println("C")
    default:
        fmt.Println("F")
    }
}
```

Порядок важен: сначала нужно ставить более конкретные или более строгие условия.

---

## 4.4. `fallthrough`

В Go следующий `case` не выполняется автоматически. Для явного проваливания есть `fallthrough`.

```go
package main

import "fmt"

func main() {
    n := 1

    switch n {
    case 1:
        fmt.Println("one")
        fallthrough
    case 2:
        fmt.Println("two")
    default:
        fmt.Println("other")
    }
}
```

Вывод:

```text
one
two
```

Важно: `fallthrough` передает управление в следующий `case` без проверки его условия. Используется редко, потому что может ухудшать читаемость.

---

## 4.5. `type switch`

`type switch` применяется, когда значение имеет интерфейсный тип, а поведение зависит от фактического типа внутри интерфейса.

```go
package main

import "fmt"

func printValue(v any) {
    switch value := v.(type) {
    case nil:
        fmt.Println("nil")
    case int:
        fmt.Println("int:", value)
    case string:
        fmt.Println("string:", value)
    case bool:
        fmt.Println("bool:", value)
    default:
        fmt.Printf("unknown type %T\n", value)
    }
}

func main() {
    printValue(10)
    printValue("hello")
    printValue(nil)
}
```

Что важно:

- синтаксис `v.(type)` разрешен только внутри `type switch`;
- переменная `value` внутри каждого `case` имеет соответствующий тип;
- `case nil` позволяет обработать пустое интерфейсное значение.

---

# 5. Цикл `for`

В Go есть один оператор цикла: `for`. Он покрывает классический цикл, `while`-стиль, бесконечный цикл и обход коллекций через `range`.

## 5.1. Классический `for`

```go
package main

import "fmt"

func main() {
    for i := 0; i < 5; i++ {
        fmt.Println(i)
    }
}
```

Порядок выполнения:

1. Один раз выполняется `i := 0`.
2. Проверяется условие `i < 5`.
3. Если условие `true`, выполняется тело цикла.
4. Выполняется пост-инструкция `i++`.
5. Снова проверяется условие.
6. Если условие `false`, цикл завершается.

---

## 5.2. `for` в стиле `while`

```go
package main

import "fmt"

func main() {
    i := 0

    for i < 5 {
        fmt.Println(i)
        i++
    }
}
```

Здесь нет отдельной инициализации и пост-инструкции. Цикл выполняется, пока условие истинно.

---

## 5.3. Бесконечный цикл

```go
package main

import "fmt"

func main() {
    i := 0

    for {
        fmt.Println(i)
        i++

        if i >= 3 {
            break
        }
    }
}
```

Бесконечный цикл часто используют в серверах, воркерах, обработчиках событий. В учебных примерах почти всегда нужен явный `break`, чтобы программа завершилась.

---

## 5.4. `range` по slice, map и string

```go
package main

import "fmt"

func main() {
    nums := []int{10, 20, 30}

    for index, value := range nums {
        fmt.Println(index, value)
    }
}
```

Если индекс не нужен:

```go
for _, value := range nums {
    fmt.Println(value)
}
```

Если нужно только первое значение:

```go
for index := range nums {
    fmt.Println(index)
}
```

Для строки `range` идет по рунам, а индекс показывает байтовое смещение.

```go
package main

import "fmt"

func main() {
    for index, r := range "Go 世界" {
        fmt.Printf("%d %c\n", index, r)
    }
}
```

Примерный вывод:

```text
0 G
1 o
2  
3 世
6 界
```

Иероглифы занимают несколько байт, поэтому индексы идут не подряд.

---

## 5.5. `break` и `continue`

`break` завершает ближайший цикл.

```go
for i := 0; i < 10; i++ {
    if i == 3 {
        break
    }
    fmt.Println(i)
}
```

`continue` переходит к следующей итерации ближайшего цикла.

```go
for i := 1; i <= 10; i++ {
    if i%2 == 0 {
        continue
    }
    fmt.Println(i)
}
```

---

## 5.6. Подводные камни в циклах

### `defer` в цикле

`defer` выполняется при выходе из функции, а не при завершении текущей итерации.

```go
package main

import "fmt"

func main() {
    for i := 0; i < 3; i++ {
        defer fmt.Println(i)
    }

    fmt.Println("end")
}
```

Вывод:

```text
end
2
1
0
```

Если ресурс нужно закрывать на каждой итерации, часто удобнее вынести тело итерации в отдельную функцию.

### Изменение коллекции во время обхода

Изменять срез или map во время обхода нужно аккуратно. Для учебной лекции достаточно правила: если логика становится неочевидной, сначала соберите изменения отдельно, а примените их после цикла.

---

# 6. Метки

Метки позволяют управлять внешним циклом из внутреннего.

```go
package main

import "fmt"

func main() {
Outer:
    for i := 1; i <= 3; i++ {
        for j := 1; j <= 3; j++ {
            if i*j > 4 {
                break Outer
            }
            fmt.Printf("%d,%d ", i, j)
        }
    }

    fmt.Println()
}
```

Обычный `break` вышел бы только из внутреннего цикла. `break Outer` выходит из цикла, помеченного меткой `Outer`.

Пример с `continue`:

```go
package main

import "fmt"

func main() {
Loop:
    for i := 1; i <= 2; i++ {
        for j := 1; j <= 3; j++ {
            if j == 2 {
                continue Loop
            }
            fmt.Printf("%d,%d ", i, j)
        }
    }

    fmt.Println()
}
```

Вывод:

```text
1,1 2,1
```

`continue Loop` переходит к следующей итерации внешнего цикла.

Практический совет: если метка делает код сложным, часто лучше вынести вложенный цикл в отдельную функцию и использовать `return`.

---

# 7. Отладка управляющих конструкций

## 7.1. Быстрый запуск

```bash
go run main.go
```

Для проверки компиляции без запуска:

```bash
go build ./...
```

Для поиска гонок данных:

```bash
go run -race main.go
```

---

## 7.2. Отладка через Delve

Пример программы:

```go
package main

import "fmt"

func main() {
    sum := 0

    for i := 1; i <= 5; i++ {
        sum += i
    }

    fmt.Println(sum)
}
```

Сборка без оптимизаций:

```bash
go build -gcflags="all=-N -l" -o sum_debug main.go
dlv exec ./sum_debug
```

Основные команды Delve:

| Команда | Значение |
|---------|----------|
| `break main.main` | поставить точку останова |
| `continue` | продолжить выполнение |
| `next` | выполнить следующую строку, не заходя внутрь функции |
| `step` | зайти внутрь функции |
| `print sum` | вывести значение переменной |
| `locals` | показать локальные переменные |
| `quit` | выйти |

Что показать на занятии:

- как меняется `i` в цикле;
- когда выполняется `i++`;
- как срабатывает `break`;
- почему правая часть `&&` или `||` иногда не вызывается.

---

# 8. Практические задания на занятии

## Задание 1. Исправить условие

Дан код:

```go
age := 16
hasID := true

if age >= 18 || hasID {
    fmt.Println("allow")
}
```

Вопрос: почему условие опасное?  
Ожидаемый ответ: оно пустит несовершеннолетнего пользователя, если у него есть документ. Правильнее:

```go
if age >= 18 && hasID {
    fmt.Println("allow")
}
```

---

## Задание 2. Безопасно проверить вложенное поле

Дан код:

```go
type Profile struct {
    Email string
}

type User struct {
    Profile *Profile
}

var user *User
```

Нужно написать условие, которое печатает email только если `user`, `Profile` и `Email` доступны.

Ожидаемое решение:

```go
if user != nil && user.Profile != nil && user.Profile.Email != "" {
    fmt.Println(user.Profile.Email)
}
```

---

## Задание 3. Предсказать порядок вывода

```go
package main

import "fmt"

func f(name string, result bool) bool {
    fmt.Println(name)
    return result
}

func main() {
    if f("A", false) && f("B", true) || f("C", true) {
        fmt.Println("ok")
    }
}
```

Ожидаемый вывод:

```text
A
C
ok
```

Пояснение:

- `A` возвращает `false`;
- для `&&` правая часть `B` не вычисляется;
- затем вычисляется правая часть `||`, то есть `C`;
- `C` возвращает `true`, поэтому условие истинно.

---

## Задание 4. Переписать вложенный цикл

Даны два среза чисел. Нужно найти первый общий элемент.

```go
a := []int{1, 2, 3}
b := []int{4, 5, 3}
```

Вариант с меткой:

```go
found := false

Search:
for _, x := range a {
    for _, y := range b {
        if x == y {
            fmt.Println("found:", x)
            found = true
            break Search
        }
    }
}

if !found {
    fmt.Println("not found")
}
```

Вариант через функцию:

```go
func findCommon(a, b []int) (int, bool) {
    for _, x := range a {
        for _, y := range b {
            if x == y {
                return x, true
            }
        }
    }

    return 0, false
}
```

Обсуждение: вариант через функцию обычно проще читать и тестировать.

---

# 9. Частые ошибки и подводные камни

Этот раздел удобно пройти после основной части лекции: он собирает ситуации, которые часто выглядят очевидно, но дают неожиданный результат.

## 9.1. `else` должен быть на той же строке

В Go автоматически вставляются точки с запятой. Поэтому `else` нельзя переносить на новую строку после закрывающей скобки `if`.

Так нельзя:

```go
if x > 0 {
    fmt.Println("positive")
}
else {
    fmt.Println("non-positive")
}
```

Так правильно:

```go
if x > 0 {
    fmt.Println("positive")
} else {
    fmt.Println("non-positive")
}
```

Идея: закрывающая `}` и `else` должны быть частью одной конструкции `if/else`, поэтому `else` пишется сразу после `}`.

---

## 9.2. Порядок `case` в `switch` без выражения

В `switch` выполняется первый подходящий `case`. Остальные варианты уже не проверяются.

```go
package main

import "fmt"

func main() {
    score := 95

    switch {
    case score >= 60:
        fmt.Println("passed")
    case score >= 90:
        fmt.Println("excellent")
    default:
        fmt.Println("failed")
    }
}
```

Вывод:

```text
passed
```

Хотя `score >= 90` тоже истинно, программа до этого `case` не дойдет. Более строгие условия обычно ставят выше:

```go
switch {
case score >= 90:
    fmt.Println("excellent")
case score >= 60:
    fmt.Println("passed")
default:
    fmt.Println("failed")
}
```

---

## 9.3. `fallthrough` не проверяет следующий `case`

`fallthrough` просто передает управление в следующий блок. Условие следующего `case` не проверяется.

```go
package main

import "fmt"

func main() {
    n := 1

    switch n {
    case 1:
        fmt.Println("one")
        fallthrough
    case 100:
        fmt.Println("one hundred")
    default:
        fmt.Println("other")
    }
}
```

Вывод:

```text
one
one hundred
```

Это может удивить: `n` не равен `100`, но блок `case 100` все равно выполнился из-за `fallthrough`.

---

## 9.4. `break` выходит из ближайшего `for`, `switch` или `select`

Если внутри цикла находится `switch`, обычный `break` выйдет только из `switch`, а не из цикла.

```go
package main

import "fmt"

func main() {
    command := "exit"

    for i := 0; i < 3; i++ {
        switch command {
        case "exit":
            fmt.Println("break switch")
            break
        }

        fmt.Println("after switch")
    }
}
```

Вывод:

```text
break switch
after switch
break switch
after switch
break switch
after switch
```

Чтобы выйти из внешнего цикла, можно использовать метку:

```go
Loop:
for {
    switch command {
    case "exit":
        break Loop
    }
}
```

Но если код находится внутри функции, часто проще и понятнее использовать `return`.

---

## 9.5. `range` по map не гарантирует порядок

Порядок обхода map в Go не определен. Нельзя писать код, который зависит от порядка ключей.

```go
package main

import "fmt"

func main() {
    users := map[string]int{
        "alice": 10,
        "bob":   20,
        "kate":  30,
    }

    for name, score := range users {
        fmt.Println(name, score)
    }
}
```

Вывод может отличаться между запусками. Если нужен стабильный порядок, сначала собирают ключи, сортируют их и затем читают значения из map.

```go
package main

import (
    "fmt"
    "sort"
)

func main() {
    users := map[string]int{
        "alice": 10,
        "bob":   20,
        "kate":  30,
    }

    names := make([]string, 0, len(users))
    for name := range users {
        names = append(names, name)
    }

    sort.Strings(names)

    for _, name := range names {
        fmt.Println(name, users[name])
    }
}
```

---

## 9.6. Изменение slice во время `range`

`range` по slice фиксирует длину обхода в начале цикла. Если внутри цикла делать `append`, новые элементы не станут автоматически частью текущего обхода.

```go
package main

import "fmt"

func main() {
    nums := []int{1, 2, 3}

    for _, n := range nums {
        fmt.Println(n)
        nums = append(nums, n*10)
    }

    fmt.Println("after:", nums)
}
```

Вывод:

```text
1
2
3
after: [1 2 3 10 20 30]
```

Цикл прошел только по исходным трем элементам. Это не ошибка компиляции, но часто ошибка логики.

---

## 9.7. `nil` interface и typed nil

Это тема ближе к интерфейсам, но полезно показать ее как предупреждение.

```go
package main

import "fmt"

type User struct {
    Name string
}

func main() {
    var user *User = nil
    var value any = user

    fmt.Println(user == nil)
    fmt.Println(value == nil)
}
```

Вывод:

```text
true
false
```

Почему так: интерфейсное значение хранит не только само значение, но и его динамический тип. В `value` лежит тип `*User` и значение `nil`, поэтому сам интерфейс уже не равен `nil`.

Это важно помнить при проверках:

```go
if value != nil {
    fmt.Println("interface is not nil")
}
```

Такой код не означает, что внутри интерфейса точно лежит безопасное для разыменования значение.

---

## 9.8. Слишком сложные условия лучше разбивать

Короткие замыкания полезны, но длинные условия быстро становятся нечитаемыми.

```go
if user != nil && user.Profile != nil && user.Profile.Email != "" && (isAdmin || user.IsActive && user.HasPaid) {
    fmt.Println("allowed")
}
```

Такой код лучше разделить:

```go
hasUser := user != nil
hasEmail := hasUser && user.Profile != nil && user.Profile.Email != ""
hasAccess := isAdmin || hasUser && user.IsActive && user.HasPaid

if hasEmail && hasAccess {
    fmt.Println("allowed")
}
```

При этом нужно следить, чтобы разбиение не сломало защиту от `nil`. Если `user` может быть `nil`, то каждое обращение к полям `user` должно оставаться за проверкой `hasUser`.

---

# 10. Итоги

Ключевые мысли лекции:

- `if` и `switch` управляют выбором ветки выполнения;
- условие в Go всегда должно быть `bool`;
- короткая инструкция в `if` и `switch` ограничивает область видимости временных переменных;
- `&&` и `||` вычисляются слева направо и используют короткое замыкание;
- короткие замыкания позволяют безопасно проверять `nil` перед опасной операцией;
- аргументы функций вычисляются до входа в функцию;
- `for` является единственным циклом в Go и покрывает несколько привычных форм;
- `break`, `continue` и метки управляют выходом из циклов;
- отладчик помогает увидеть реальный порядок выполнения программы.

---

# 11. Домашнее задание

1. Написать FizzBuzz двумя способами: через `if/else` и через `switch`.
2. Написать функцию `CanAccess(age int, hasID bool, isAdmin bool) bool`.
3. Написать безопасную функцию `UserEmail(user *User) (string, bool)`.
4. Реализовать поиск первого общего элемента в двух срезах.
5. Запустить один из примеров под `dlv` и посмотреть значения переменных в цикле.

---

# 12. Дополнительные материалы

- [Go Tour: Flow control](https://go.dev/tour/flowcontrol/1)
- [Effective Go: Control structures](https://go.dev/doc/effective_go#if)
- [Go Spec: Order of evaluation](https://go.dev/ref/spec#Order_of_evaluation)
- [Go Spec: For statements](https://go.dev/ref/spec#For_statements)
- [Go Spec: Switch statements](https://go.dev/ref/spec#Switch_statements)
- [Delve Debugger](https://github.com/go-delve/delve)
