.PHONY: dev down logs ps restart

dev:
	docker compose stop
	docker compose rm -f
	docker compose up --build

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker compose ps

restart:
	make down && docker compose up --build -d
