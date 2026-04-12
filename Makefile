.PHONY: run build db-up db-down migrate-up migrate-down

include .env
export

run:
	go run cmd/bot/main.go

build:
	go build -o bin/bot cmd/bot/main.go

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

migrate-up:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable" down

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)