# axle

Axle orchestrates AI agents through your entire product lifecycle — so you focus on decisions, not busywork.

## Architecture

| Service      | Port | Responsibility                                                   |
| :----------- | :--- | :--------------------------------------------------------------- |
| **BFF**      | 8080 | REST/ConnectRPC API for frontend, CRUD, business logic           |
| **Gateway**  | 8081 | Real-time events, SSE streaming, agent pipeline orchestration    |
| **LLM**      | 8082 | Bifrost wrapper (LLM provider), RAG pipeline, AI agent execution |
| **Frontend** | 5173 | Vue 3 SPA                                                        |

All services communicate via **NATS JetStream** (except Frontend ↔ BFF and Frontend ↔ Gateway).

## CLI Tools

### 1. contracts (Protobuf, ConnectRPC)

```bash
cd contracts/generated
pnpm install
cd ..
buf generate     																		# buf generate → generates Go + TypeScript code
buf lint        																		# buf lint
buf breaking --against '.git#branch=main'     			# buf breaking (check backwards compatibility)
```

**Generated code:**

- `contracts/generated/go/` — Go (protoc-gen-go + protoc-gen-connect-go)
- `contracts/generated/es/` — TypeScript (protoc-gen-es v2, Connect v2)

### 2. frontend (Vue 3, Vite, Pinia, Naive UI, vue-i18n)

```bash
cd frontend
pnpm install
pnpm dev          # dev server → http://localhost:5173
pnpm build        # production build
pnpm type-check   # TypeScript check
pnpm format       # prettier
pnpm lint         # oxlint + eslint
```

### 3. services/bff (Go)

```bash
cd services/bff
go mod tidy
go build ./...
go run ./cmd/bff          # starts on :8080
```

**DB migrations (goose):**

```bash
go run github.com/pressly/goose/v3/cmd/goose@latest \
  -dir db/migrations postgres "$POSTGRES_DSN" up
```

**sqlc (generate store code):**

```bash
sqlc generate    # requires sqlc CLI: https://sqlc.dev
```

### 4. services/gateway (Go)

```bash
cd services/gateway
go mod tidy
go run ./cmd/gateway      # starts on :8081
```

### 5. services/llm (Go)

```bash
cd services/llm
go mod tidy
go run ./cmd/llm        # starts on :8082
```

## Dev environment (Docker Compose)

```bash
cd deploy
cp .env.example .env
docker compose -f docker-compose.dev.yaml up -d
```

Services started:

- **PostgreSQL 16** with `pgvector` → localhost:5432
- **Redis 7** → localhost:6379
- **NATS** with JetStream → localhost:4222, monitor: localhost:8222
- **MinIO** (S3) → localhost:9000, console: localhost:9001
