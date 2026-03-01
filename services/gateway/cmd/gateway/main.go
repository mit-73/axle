package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/ApeironFoundation/axle/contracts/go/gateway/v1/gen_gateway_v1connect"
	"github.com/ApeironFoundation/axle/gateway/internal/config"
	"github.com/ApeironFoundation/axle/gateway/internal/enterprise"
	"github.com/ApeironFoundation/axle/gateway/internal/health"
	"github.com/ApeironFoundation/axle/gateway/internal/hub"
	"github.com/ApeironFoundation/axle/gateway/internal/natsclient"
	"github.com/ApeironFoundation/axle/gateway/internal/streaming"
)

func main() {
	// ── Structured logging ───────────────────────────────────────────────────
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// ── Config ───────────────────────────────────────────────────────────────
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to load config")
	}

	lvl, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		lvl = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(lvl)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// ── Redis ────────────────────────────────────────────────────────────────
	log.Info().Str("url", cfg.RedisURL).Msg("connecting to redis")
	redisOpts, err := redis.ParseURL(cfg.RedisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("invalid redis URL")
	}
	rdb := redis.NewClient(redisOpts)
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatal().Err(err).Msg("redis ping failed")
	}
	defer func() { _ = rdb.Close() }()
	log.Info().Msg("redis connected")

	// ── NATS ─────────────────────────────────────────────────────────────────
	log.Info().Str("url", cfg.NatsURL).Msg("connecting to nats")
	natsConns, err := natsclient.Connect(ctx, cfg.NatsURL)
	if err != nil {
		log.Fatal().Err(err).Msg("nats connect failed")
	}
	defer natsConns.NC.Drain() //nolint:errcheck
	log.Info().Msg("nats connected")

	// ── Enterprise registry ──────────────────────────────────────────────────
	_ = enterprise.NewRegistry()

	// ── Streaming hub ────────────────────────────────────────────────────────
	eventHub := hub.New()

	// Subscribe to NATS events topic and fan-out to hub
	_, err = natsConns.NC.Subscribe("axle.events.>", func(msg *nats.Msg) {
		if msg.Subject == "axle.events.test.ping" {
			log.Info().Str("subject", msg.Subject).Int("bytes", len(msg.Data)).Msg("dev-only ping event received by gateway")
		}
		eventHub.Publish(context.Background(), msg.Data)
	})
	if err != nil {
		log.Fatal().Err(err).Msg("nats subscribe failed")
	}

	// ── Health checker ───────────────────────────────────────────────────────
	checker := health.NewChecker(rdb, natsConns.NC)

	// ── Router ───────────────────────────────────────────────────────────────
	r := chi.NewRouter()

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			id := req.Header.Get("X-Request-Id")
			if id == "" {
				id = uuid.New().String()
			}
			w.Header().Set("X-Request-Id", id)
			next.ServeHTTP(w, req)
		})
	})

	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{
			"Authorization", "Content-Type", "Connect-Protocol-Version",
			"Connect-Timeout-Ms", "X-Request-Id",
		},
	}).Handler)

	r.Get("/health", checker.HealthHandler)
	r.Get("/ready", checker.ReadyHandler)

	// ConnectRPC streaming service
	connectMux := http.NewServeMux()
	connectMux.Handle(gen_gateway_v1connect.NewStreamingServiceHandler(
		streaming.NewHandler(eventHub),
	))
	r.Handle("/gateway.v1.*", connectMux)

	// ── HTTP server ──────────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           h2c.NewHandler(r, &http2.Server{}),
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      0, // streaming — no write timeout
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Info().Str("addr", addr).Msg("gateway listening")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatal().Err(err).Msg("server error")
		}
	}()

	// ── Graceful shutdown ────────────────────────────────────────────────────
	<-ctx.Done()
	log.Info().Msg("shutting down…")

	shutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(shutCtx); err != nil {
		log.Error().Err(err).Msg("graceful shutdown failed")
	}
	log.Info().Msg("gateway stopped")
}
