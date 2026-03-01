# axle

Axle orchestrates AI agents through your entire product lifecycle — so you focus on decisions, not busywork.

## Architecture

| Service      | Port | Responsibility                                                   |
| :----------- | :--- | :--------------------------------------------------------------- |
| **BFF**      | 9001 | REST/ConnectRPC API for frontend, CRUD, business logic           |
| **Gateway**  | 9002 | Real-time events, SSE streaming, agent pipeline orchestration    |
| **LLM**      | 9003 | Bifrost wrapper (LLM provider), RAG pipeline, AI agent execution |
| **Frontend** | 5173 | Vue 3 SPA                                                        |

All services communicate via **NATS JetStream** (except Frontend ↔ BFF and Frontend ↔ Gateway).

## Prerequisites

Install these tools before running the project:

- **just** (task runner used across this repo)
- **Go** (for backend services and Go-based generators)
- **Node.js + pnpm** (for frontend and generated TS contracts)
- **buf** (Protobuf lint/breaking/codegen)
- **Docker + Docker Compose** (local infra: Postgres, Redis, NATS, MinIO)

You can verify your setup with:

```bash
just init
```

## Quick start

```bash
just init
just dev-up
```

Then run services in separate terminals:

```bash
just pipeline-debug-help
```

## Core workflows

This repository standardizes workflows via `justfile`.
Prefer `just ...` commands instead of ad-hoc service-specific commands.

Common commands:

- `just dev-up` / `just dev-down` / `just dev-logs`
- `just contracts-generate`
- `just db-generate`
- `just pipeline-validate`

See all available commands:

```bash
just --list
```

## Dev environment (Docker Compose)

Services started:

- **PostgreSQL 16** with `pgvector` → localhost:5432
- **Redis 7** → localhost:6379
- **NATS** with JetStream → localhost:4222, monitor: localhost:8222
- **MinIO** (S3) → localhost:9000, console: localhost:9001
