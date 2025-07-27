.PHONY: dev down logs ps restart

start:
	docker compose -f docker-compose.yml -f local-compose-override.yml up --build -d

dev:
	docker compose stop
	docker compose rm -f
	make start

down:
	docker compose down

logs:
	docker compose logs -f

ps:
	docker compose ps

restart:
	make down && make start
