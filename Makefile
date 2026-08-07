MIGRATE ?= C:\Users\blanc\go\bin\migrate.exe

DB_HOST ?= localhost
DB_PORT ?= 5432
DB_NAME ?= neurosis
DB_USER ?= postgres
DB_PASSWORD ?= postgres
DB_SSLMODE ?= disable

DB_URL := postgres://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSLMODE)
MIGRATIONS_PATH := ./migrations

.PHONY: migrate-up migrate-down migrate-force run

migrate-up:
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" up

migrate-down:
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" down

migrate-force:
ifndef VERSION
	$(error VERSION is required. Use: make migrate-force VERSION=<n>)
endif
	$(MIGRATE) -path $(MIGRATIONS_PATH) -database "$(DB_URL)" force $(VERSION)

run:
	DB_HOST=$(DB_HOST) DB_PORT=$(DB_PORT) DB_NAME=$(DB_NAME) DB_USER=$(DB_USER) DB_PASSWORD=$(DB_PASSWORD) DB_SSLMODE=$(DB_SSLMODE) go run ./cmd/api/main.go
