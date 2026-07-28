# 15. Конфигурация - примеры

Учебные проекты к теме конфигурации Go-приложений (курс, Go 1.26.x).

Подробный разбор: [lecture.md](./lecture.md).

## Каталоги

| Каталог | Содержание | README |
|---------|------------|--------|
| `json/` | `encoding/json` + `config.json` | [README](./json/README.md) |
| `yaml/` | `gopkg.in/yaml.v3` + `config.yaml` | [README](./yaml/README.md) |
| `ini/` | `gopkg.in/ini.v1` + `config.ini` | [README](./ini/README.md) |
| `env/` | `godotenv` + `os.Getenv` + `.env` | [README](./env/README.md) |
| `cli/flag-cli/` | стандартный пакет `flag` | [README](./cli/flag-cli/README.md) |
| `cli/cobra-cli/` | подкоманды Cobra | [README](./cli/cobra-cli/README.md) |
| `viper/` | Viper + Cobra: файл и флаг | [README](./viper/README.md) |
| `all/` | defaults + file + env + flags | [README](./all/README.md) |

Запуск типичного примера:

```bash
cd json   # или yaml, ini, env, ...
go mod tidy   # если есть go.mod
go run .
```

Для `cli/flag-cli` модуля может не быть - достаточно `go run .` из каталога примера.
