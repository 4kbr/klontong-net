## Kontrak OpenAPI Klontong Net. Bungkus tipis pnpm scripts.
.DEFAULT_GOAL := help

.PHONY: install lint bundle mock preview ci help

install: ## pasang tooling (redocly, spectral, prism)
	pnpm install

lint: ## jalankan Spectral lint
	pnpm run lint

bundle: ## rakit jadi dist/openapi.yaml (satu file)
	pnpm run bundle

mock: ## jalankan Prism mock server di :4010
	pnpm run mock

preview: ## pratinjau dokumentasi Redoc di browser
	pnpm run preview

ci: ## gate: lint + bundle
	pnpm run ci

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN{FS=":.*?## "}{printf "  \033[36m%-10s\033[0m %s\n", $$1, $$2}'
