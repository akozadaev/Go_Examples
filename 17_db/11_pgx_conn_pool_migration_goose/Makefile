.PHONY: help install up status version down redo reset up-to down-to

GOOSE ?= goose
GOOSE_VERSION ?= v3.27.3
MIGRATIONS_DIR ?= migrations
DB_URL ?= postgres://ibs:ibs@localhost:5432/ibs_goose?sslmode=disable

GOOSE_CMD = $(GOOSE) -dir $(MIGRATIONS_DIR) postgres "$(DB_URL)"

help:
	@echo "Команды управления миграциями Goose:"
	@echo "  make install             Установить Goose $(GOOSE_VERSION)"
	@echo "  make up                  Применить все новые миграции"
	@echo "  make status              Показать состояние миграций"
	@echo "  make version             Показать текущую версию схемы"
	@echo "  make down                Откатить последнюю миграцию"
	@echo "  make redo                Откатить и повторно применить последнюю миграцию"
	@echo "  make reset               Откатить все миграции"
	@echo "  make up-to VERSION=2     Применить миграции до заданной версии"
	@echo "  make down-to VERSION=1   Откатить миграции до заданной версии"
	@echo ""
	@echo "Параметры можно переопределить:"
	@echo "  make status DB_URL='postgres://user:password@host:5432/database?sslmode=disable'"

install:
	go install github.com/pressly/goose/v3/cmd/goose@$(GOOSE_VERSION)

up:
	$(GOOSE_CMD) up

status:
	$(GOOSE_CMD) status

version:
	$(GOOSE_CMD) version

down:
	$(GOOSE_CMD) down

redo:
	$(GOOSE_CMD) redo

reset:
	$(GOOSE_CMD) reset

up-to:
	@test -n "$(VERSION)" || (echo "Укажите версию: make up-to VERSION=2"; exit 1)
	$(GOOSE_CMD) up-to $(VERSION)

down-to:
	@test -n "$(VERSION)" || (echo "Укажите версию: make down-to VERSION=1"; exit 1)
	$(GOOSE_CMD) down-to $(VERSION)
