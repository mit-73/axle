# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY contracts/generated /src/contracts/generated
COPY services/gateway /src/services/gateway

WORKDIR /src/services/gateway
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /bin/gateway \
    ./cmd/gateway

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12 AS runtime

COPY --from=builder /bin/gateway /gateway

EXPOSE 9002

USER nonroot:nonroot

ENTRYPOINT ["/gateway"]
