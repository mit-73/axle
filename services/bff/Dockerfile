# ── Build stage ───────────────────────────────────────────────────────────────
FROM golang:1.26-alpine AS builder

WORKDIR /src

COPY contracts/generated /src/contracts/generated
COPY db /src/db
COPY services/bff /src/services/bff

WORKDIR /src/services/bff
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \
    -o /bin/bff \
    ./cmd/bff

# ── Runtime stage ─────────────────────────────────────────────────────────────
FROM gcr.io/distroless/static-debian12 AS runtime

COPY --from=builder /bin/bff /bff

EXPOSE 9001

USER nonroot:nonroot

ENTRYPOINT ["/bff"]
