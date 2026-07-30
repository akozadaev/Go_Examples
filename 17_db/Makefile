.PHONY: help up down ps

help:
	@echo "Доступные команды:"
	@echo "  make up    Запустить PostgreSQL в Docker"
	@echo "  make down  Остановить PostgreSQL и удалить учебные данные"
	@echo "  make ps    Показать состояние контейнера"

up:
	docker compose up -d

down:
	docker compose down -v

ps:
	docker compose ps
