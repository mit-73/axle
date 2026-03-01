# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY contracts/generated /src/contracts/generated
COPY db /src/db
COPY services/llm /src/services/llm

WORKDIR /src/services/llm
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /bin/llm \
    ./cmd/llm

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12 AS runtime

COPY --from=builder /bin/llm /llm

EXPOSE 9003

USER nonroot:nonroot

ENTRYPOINT ["/llm"]
