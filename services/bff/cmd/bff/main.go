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
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/ApeironFoundation/axle/bff/internal/config"
	"github.com/ApeironFoundation/axle/bff/internal/db"
	"github.com/ApeironFoundation/axle/bff/internal/enterprise"
	"github.com/ApeironFoundation/axle/bff/internal/handler"
	"github.com/ApeironFoundation/axle/bff/internal/health"
	"github.com/ApeironFoundation/axle/bff/internal/middleware"
	"github.com/ApeironFoundation/axle/bff/internal/natsclient"
	"github.com/ApeironFoundation/axle/bff/internal/redisclient"

	"github.com/ApeironFoundation/axle/contracts/generated/go/bff/v1/bffv1connect"
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

	// ── PostgreSQL ───────────────────────────────────────────────────────────
	log.Info().Str("dsn_host", cfg.PostgresDSN).Msg("connecting to postgres")
	pool, err := db.Connect(ctx, cfg.PostgresDSN)
	if err != nil {
		log.Fatal().Err(err).Msg("postgres connect failed")
	}
	defer pool.Close()
	log.Info().Msg("postgres connected")

	// ── Redis ────────────────────────────────────────────────────────────────
	log.Info().Str("url", cfg.RedisURL).Msg("connecting to redis")
	rdb, err := redisclient.Connect(ctx, cfg.RedisURL)
	if err != nil {
		log.Fatal().Err(err).Msg("redis connect failed")
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

	// ── Enterprise registry (OSS no-op) ──────────────────────────────────────
	_ = enterprise.NewRegistry()

	// ── Health checker ───────────────────────────────────────────────────────
	checker := health.NewChecker(pool, rdb, natsConns.NC)

	// ── Router ───────────────────────────────────────────────────────────────
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Recovery)
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Authorization", "Content-Type", "Connect-Protocol-Version",
			"Connect-Timeout-Ms", "Grpc-Timeout", "X-Request-Id",
		},
		ExposedHeaders: []string{"Grpc-Status", "Grpc-Message", "Connect-Status"},
	}).Handler)

	// Infrastructure endpoints (no auth)
	r.Get("/health", checker.HealthHandler)
	r.Get("/ready", checker.ReadyHandler)

	// ConnectRPC handlers (served over HTTP/2 h2c, no auth middleware for smoke test)
	connectMux := http.NewServeMux()
	connectMux.Handle(bffv1connect.NewProjectServiceHandler(
		&handler.ProjectsHandler{Pool: pool},
	))
	connectMux.Handle(bffv1connect.NewUserServiceHandler(
		&handler.UsersHandler{},
	))

	// Route all ConnectRPC traffic: /bff.v1.<Service>/<Method>
	r.HandleFunc("/bff.v1.*", connectMux.ServeHTTP)

	// API routes (auth middleware applied)
	r.Group(func(r chi.Router) {
		r.Use(middleware.Auth)
		// TODO: authenticated REST endpoints
	})

	// ── HTTP server ──────────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           h2c.NewHandler(r, &http2.Server{}),
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	// Start server in background
	go func() {
		log.Info().Str("addr", addr).Msg("bff listening")
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
	log.Info().Msg("bff stopped")
}
