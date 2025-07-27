include .env

.PHONY: build up down logs
build:
	export COMPOSE_BAKE=true && \
	docker-compose build

up:
	docker-compose up -d

down:
	docker-compose down -v

logs:
	docker-compose logs -f api

.PHONY: goose-add goose-up goose-down goose-status

goose-add:
	goose -dir ./migrations postgres "$(DATABASE_DSN)" create $(NAME) sql

goose-up:
	goose -dir ./migrations postgres "$(DATABASE_DSN)" up

goose-down:
	goose -dir ./migrations postgres "$(DATABASE_DSN)" down

goose-status:
	goose -dir ./migrations postgres "$(DATABASE_DSN)" status

.PHONY: lint
lint:
	golangci-lint run