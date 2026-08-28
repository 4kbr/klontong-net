include .env
export

APP_NAME       := klontong-api
BIN_DIR        := bin
MIGRATIONS_DIR := migrations
PKG            := ./...

.DEFAULT_GOAL := help

.PHONY: dev
dev: ## Jalankan API dengan hot reload
	air -c .air.toml

.PHONY: run
run: ## Jalankan API sekali
	go run ./cmd/api

.PHONY: worker
worker: ## Jalankan background worker
	go run ./cmd/worker

.PHONY: build
build: ## Compile kedua binary
	go build -o $(BIN_DIR)/$(APP_NAME) ./cmd/api
	go build -o $(BIN_DIR)/$(APP_NAME)-worker ./cmd/worker

.PHONY: migrate-create
migrate-create: ## Buat migrasi baru. Pakai: make migrate-create name=xxx
	@test -n "$(name)" || (echo "ERROR: name= wajib diisi"; exit 1)
	migrate create -ext sql -dir $(MIGRATIONS_DIR) -seq $(name)

.PHONY: migrate-up
migrate-up: ## Terapkan migrasi
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" up

.PHONY: migrate-down
migrate-down: ## Rollback 1 migrasi
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" down 1

.PHONY: migrate-force
migrate-force: ## Paksa versi saat dirty. Pakai: make migrate-force version=3
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" force $(version)

.PHONY: migrate-version
migrate-version: ## Versi migrasi aktif
	migrate -path $(MIGRATIONS_DIR) -database "$(DATABASE_URL)" version

.PHONY: seed
seed: ## Isi data contoh
	go run ./cmd/seed

.PHONY: test
test: ## Unit test
	go test -race -short $(PKG)

.PHONY: test-integration
test-integration: ## Semua test (butuh docker)
	go test -race -count=1 $(PKG)

.PHONY: cover
cover: ## Laporan coverage
	go test -race -coverprofile=coverage.out $(PKG)
	go tool cover -html=coverage.out -o coverage.html

.PHONY: lint
lint: ## golangci-lint
	golangci-lint run

.PHONY: fmt
fmt: ## Format kode
	gofmt -s -w .
	goimports -w .

.PHONY: check
check: fmt lint test ## Gate sebelum commit

.PHONY: tools
tools: ## Install tool CLI
	go install github.com/air-verse/air@latest
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest

.PHONY: help
help:
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
