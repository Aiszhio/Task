include .env

.PHONY: build run down logs
build:
	export COMPOSE_BAKE=true && \
	docker-compose build

run:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f api

.PHONY: goose-up goose-down
goose-up:
	docker-compose run --rm api \
    	  goose -dir ./migrations postgres "$(DATABASE_DSN)" up

goose-down:
	docker-compose run --rm api \
    	  goose -dir ./migrations postgres "$(DATABASE_DSN)" down