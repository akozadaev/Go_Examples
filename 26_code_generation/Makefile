.PHONY: generate clean install-tools test bench all

# Установка всех необходимых инструментов
install-tools:
	@echo "Установка инструментов кодогенерации..."
	go install golang.org/x/tools/cmd/stringer@latest
	go install github.com/campoy/jsonenums@latest
	go install github.com/mailru/easyjson/easyjson@latest
	go install github.com/josharian/impl@latest
	@echo "Инструменты установлены!"

# Генерация кода для всех примеров
generate:
	@echo "Генерация кода..."
	@echo "1. Генерация stringer для Status..."
	cd stringer_example && go generate
	@echo "2. Генерация jsonenums для Priority..."
	cd jsonenums_example && go generate
	@echo "3. Генерация easyjson для User..."
	cd easyjson_example && go generate
	@echo "4. Генерация мапперов..."
	cd mapper_example && go generate
	@echo "Генерация завершена!"

# Генерация через go generate (альтернативный способ)
generate-all:
	@echo "Запуск go generate для всех пакетов..."
	go generate ./...
	@echo "Генерация завершена!"

# Запуск всех примеров
run:
	@echo "=== Запуск примеров ==="
	@echo "\n1. Stringer пример:"
	@cd cmd/stringer && go run .
	@echo "\n2. Jsonenums пример:"
	@cd cmd/jsonenums && go run .
	@echo "\n3. Easyjson пример:"
	@cd cmd/easyjson && go run .
	@echo "\n4. Impl пример:"
	@cd cmd/impl && go run .
	@echo "\n5. Mapper пример:"
	@cd cmd/mapper-demo && go run .

# Запуск тестов
test:
	go test ./...

# Запуск бенчмарков
bench:
	go test -bench=. -benchmem ./easyjson_example

# Очистка сгенерированных файлов
clean:
	@echo "Очистка сгенерированных файлов..."
	find . -name "*_string.go" -delete
	find . -name "*_jsonenums.go" -delete
	find . -name "*_easyjson.go" -delete
	find . -name "*_mapper_gen.go" -delete
	find . -name "*.log" -delete
	@echo "Очистка завершена!"

# Полная сборка и запуск
all: install-tools generate run

# Помощь
help:
	@echo "Доступные команды:"
	@echo "  make install-tools  - Установить все инструменты кодогенерации"
	@echo "  make generate       - Сгенерировать код для всех примеров"
	@echo "  make generate-all   - Запустить go generate для всех пакетов"
	@echo "  make run            - Запустить все примеры"
	@echo "  make test           - Запустить тесты"
	@echo "  make bench          - Запустить бенчмарки"
	@echo "  make clean          - Удалить сгенерированные файлы"
	@echo "  make all            - Установить инструменты, сгенерировать код и запустить примеры"
	@echo "  make help           - Показать эту справку"

