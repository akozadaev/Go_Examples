# Минипроекты к лекции

Каждая папка содержит отдельный `main.go`, который можно запускать независимо.

Из корня репозитория:

```bash
go run 2_if_else_example/examples/01_basic_if_else/main.go
```

Если запускаете из папки `2_if_else_example`:

```bash
go run examples/01_basic_if_else/main.go
```

В этом репозитории лучше запускать именно файл `main.go`, а не папку через `go run ./...`: в полном пути проекта есть кириллические символы, из-за которых Go может считать путь некорректным import path вне модуля.

Список примеров:

| Папка | Тема |
|-------|------|
| `01_basic_if_else` | базовый `if/else` |
| `02_if_with_short_statement` | `if` с короткой инструкцией |
| `03_early_return` | ранний выход и основной сценарий |
| `04_if_execution_order` | порядок выполнения `if` |
| `05_expression_evaluation_order` | порядок вычисления выражений |
| `06_function_arguments` | вычисление аргументов функций |
| `07_defer_arguments` | вычисление аргументов `defer` |
| `08_short_circuit_and` | короткое замыкание `&&` |
| `09_short_circuit_or` | короткое замыкание `||` |
| `10_short_circuit_mixed` | смешанное условие `&&` и `||` |
| `11_nil_pointer_guard` | безопасная проверка указателя на `nil` |
| `12_nested_nil_guard` | проверка вложенных полей |
| `13_nil_slice_map_channel` | особенности `nil` slice/map/channel |
| `14_classic_switch` | классический `switch` |
| `15_switch_evaluation_order` | порядок вычисления `switch` |
| `16_switch_without_expression` | `switch` без выражения |
| `17_fallthrough` | явный `fallthrough` |
| `18_type_switch` | `type switch` |
| `19_for_classic` | классический `for` |
| `20_for_while_style` | `for` в стиле `while` |
| `21_for_infinite_with_break` | бесконечный цикл с `break` |
| `22_range_slice` | `range` по срезу |
| `23_range_string` | `range` по строке |
| `24_break_continue` | `break` и `continue` |
| `25_defer_in_loop` | `defer` в цикле |
| `26_labels_break` | метка с `break` |
| `27_labels_continue` | метка с `continue` |
| `28_find_common_with_label` | поиск с меткой |
| `29_find_common_with_function` | поиск через функцию |
| `30_debug_sum` | пример для отладки |
| `31_fizzbuzz_if_else` | FizzBuzz через `if/else` |
| `32_fizzbuzz_switch` | FizzBuzz через `switch` |
| `33_if_init_condition_else` | полная форма `if init; condition { ... } else { ... }` |
