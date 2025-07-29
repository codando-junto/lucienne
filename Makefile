.PHONY: start
start:
	docker compose -f docker-compose.yml -f local-compose-override.yml up --build -d

.PHONY: dev
dev:
	docker compose stop
	docker compose rm -f
	$(MAKE) start

.PHONY: down
down:
	docker compose down

.PHONY: logs
logs:
	docker compose logs -f

.PHONY: ps
ps:
	docker compose ps

.PHONY: restart
restart: down start
