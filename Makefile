.PHONY: dev down logs ps restart

dev:
	docker compose stop
	docker compose rm -f
	docker compose -f docker-compose.yml -f local-compose-override.yml up --build

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker compose ps

restart:
	make down && docker compose up --build -d
