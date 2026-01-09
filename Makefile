include .env
export $(shell sed 's/=.*//' .env)
MGDIR = sql/spoti/migration

dev:
	go run ./cmd

up:
	goose -dir $(MGDIR) postgres ${POSTGRES_URL} up

create_%:
	goose -dir $(MGDIR) create $* sql

status:
	goose -dir $(MGDIR) postgres ${POSTGRES_URL} status

reset:
	goose -dir $(MGDIR) postgres ${POSTGRES_URL} reset
	goose -dir $(MGDIR) postgres ${POSTGRES_URL} up

down:
	goose -dir $(MGDIR) postgres ${POSTGRES_URL} down

swag:
	swag init -g cmd/main.go