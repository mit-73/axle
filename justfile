default:
	@just --list

# =====================================
# INIT COMMANDS (run once on new machine)
# =====================================

init: init-check init-deps init-done

init-check:
	@echo "Checking required CLI tools..."
	just --version
	go version
	pnpm --version
	buf --version
	docker --version
	docker compose version

init-deps: contracts-install frontend-install go-tidy-all

init-done:
	@echo "Initialization complete! You can now run 'just dev-up' to start the development environment."

contracts-install:
	(cd contracts/generated && pnpm install)

frontend-install:
	(cd frontend && pnpm install)

go-tidy-all:
	(cd db && go mod tidy)
	(cd services/bff && go mod tidy)
	(cd services/gateway && go mod tidy)
	(cd services/llm && go mod tidy)


# =====================================
# DAY-TO-DAY COMMANDS (regular usage)
# =====================================

# Contracts
contracts-generate:
	(cd contracts && buf generate)

contracts-lint:
	(cd contracts && buf lint)

contracts-breaking:
	(cd contracts && buf breaking --against '.git#branch=main')

# Database
db-generate:
	(cd db && go tool sqlc generate)

db-migrate-up:
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir db/migrations postgres "$POSTGRES_DSN" up

db-migrate-down:
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir db/migrations postgres "$POSTGRES_DSN" down

# Frontend
frontend-dev:
	(cd frontend && pnpm dev)

frontend-build:
	(cd frontend && pnpm build)

frontend-type-check:
	(cd frontend && pnpm type-check)

frontend-lint:
	(cd frontend && pnpm lint)

frontend-format:
	(cd frontend && pnpm format)

# Services
bff-run:
	(cd services/bff && go run ./cmd/bff)

bff-build:
	(cd services/bff && go build ./...)

gateway-run:
	(cd services/gateway && go run ./cmd/gateway)

gateway-build:
	(cd services/gateway && go build ./...)

llm-run:
	(cd services/llm && go run ./cmd/llm)

llm-build:
	(cd services/llm && go build ./...)


# =====================================
# PIPELINE COMMANDS (debug and validation flows)
# =====================================

dev-up:
	docker compose -f deploy/docker-compose.infra.yaml up -d

dev-down:
	docker compose -f deploy/docker-compose.infra.yaml down

dev-logs:
	docker compose -f deploy/docker-compose.infra.yaml logs -f

pipeline-prepare: dev-up contracts-generate db-generate

pipeline-validate: contracts-lint frontend-type-check frontend-lint bff-build gateway-build llm-build

pipeline-debug-help:
	@echo "Pipeline debug flow:"
	@echo "  1) just dev-up"
	@echo "  2) just bff-run       (terminal #1)"
	@echo "  3) just gateway-run   (terminal #2)"
	@echo "  4) just llm-run       (terminal #3)"
	@echo "  5) just frontend-dev  (terminal #4)"
