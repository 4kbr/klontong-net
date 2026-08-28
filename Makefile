BACKEND := apps/backend

.DEFAULT_GOAL := help

.PHONY: setup
setup: ## Setup awal dari nol
	cp -n .env.example .env || true
	cp -n $(BACKEND)/.env.example $(BACKEND)/.env || true
	$(MAKE) -C $(BACKEND) tools
	$(MAKE) up
	@echo ""
	@echo "Berikutnya: isi kredensial payment gateway di apps/backend/.env, lalu 'make migrate' dan 'make dev'."

.PHONY: up
up: ## Nyalakan infra dev
	docker compose -f docker-compose.dev.yml up -d

.PHONY: down
down: ## Matikan infra dev
	docker compose -f docker-compose.dev.yml down

.PHONY: down-v
down-v: ## Matikan infra dev + hapus volume (RESET DATABASE)
	docker compose -f docker-compose.dev.yml down -v

.PHONY: logs
logs: ## Ikuti log infra dev
	docker compose -f docker-compose.dev.yml logs -f

.PHONY: dev
dev: ## Jalankan API dengan hot reload
	$(MAKE) -C $(BACKEND) dev

.PHONY: worker
worker: ## Jalankan background worker
	$(MAKE) -C $(BACKEND) worker

.PHONY: tunnel
tunnel: ## Tunnel publik agar payment gateway bisa mengirim webhook
	@echo "Payment gateway harus bisa memanggil endpoint webhook dari internet."
	@echo "Jalankan salah satu, lalu daftarkan URL-nya di dashboard gateway:"
	@echo "  ngrok http 8080"
	@echo "  cloudflared tunnel --url http://localhost:8080"

.PHONY: migrate
migrate: ## Terapkan migrasi
	$(MAKE) -C $(BACKEND) migrate-up

.PHONY: migrate-create
migrate-create: ## Buat migrasi baru. Pakai: make migrate-create name=xxx
	$(MAKE) -C $(BACKEND) migrate-create name=$(name)

.PHONY: seed
seed: ## Isi data contoh
	$(MAKE) -C $(BACKEND) seed

.PHONY: test
test: ## Unit test
	$(MAKE) -C $(BACKEND) test

.PHONY: test-integration
test-integration: ## Semua test (butuh docker)
	$(MAKE) -C $(BACKEND) test-integration

.PHONY: lint
lint: ## Lint
	$(MAKE) -C $(BACKEND) lint

.PHONY: check
check: ## Gate sebelum commit
	$(MAKE) -C $(BACKEND) check

.PHONY: help
help: ## Tampilkan bantuan ini
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-18s\033[0m %s\n", $$1, $$2}'
