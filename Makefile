include .env
export $(shell sed 's/=.*//' .env)

PostgresDIR = sql/migrations/postgres
ClickhouseDIR = sql/migrations/clickhouse

test:
	go test ./...

test-cache:
	go clean -testcache && go test ./...

dev:
	go run ./cmd

up:
	goose -dir $(PostgresDIR) postgres ${POSTGRES_URL} up

create_%:
	goose -dir $(PostgresDIR) create $* sql

status:
	goose -dir $(PostgresDIR) postgres ${POSTGRES_URL} status

reset:
	goose -dir $(PostgresDIR) postgres ${POSTGRES_URL} reset
	goose -dir $(PostgresDIR) postgres ${POSTGRES_URL} up

down:
	goose -dir $(PostgresDIR) postgres ${POSTGRES_URL} down

# CLICKHOUSE
ch_up:
	goose clickhouse -dir $(ClickhouseDIR) ${CLICKHOUSEURL} up

ch_create_%:
	goose -dir $(ClickhouseDIR) create $* sql

ch_status:
	goose clickhouse -dir $(ClickhouseDIR) ${CLICKHOUSEURL} status

ch_reset:
	goose clickhouse -dir $(ClickhouseDIR) ${CLICKHOUSEURL} reset
	goose clickhouse -dir $(ClickhouseDIR) ${CLICKHOUSEURL} up

ch_down:
	goose clickhouse -dir $(ClickhouseDIR) ${CLICKHOUSEURL} down

seed_db:
	go run sql/migrations/seed/*.go

swag:
	swag init -g cmd/main.go