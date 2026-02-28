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

	"github.com/ApeironFoundation/axle/llm/internal/bifrost"
	"github.com/ApeironFoundation/axle/llm/internal/config"
	"github.com/ApeironFoundation/axle/llm/internal/db"
	"github.com/ApeironFoundation/axle/llm/internal/enterprise"
	"github.com/ApeironFoundation/axle/llm/internal/handler"
	"github.com/ApeironFoundation/axle/llm/internal/health"
	"github.com/ApeironFoundation/axle/llm/internal/natsclient"

	"github.com/ApeironFoundation/axle/contracts/generated/go/ai/v1/aiv1connect"
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

	// ── NATS ─────────────────────────────────────────────────────────────────
	log.Info().Str("url", cfg.NatsURL).Msg("connecting to nats")
	natsConns, err := natsclient.Connect(ctx, cfg.NatsURL)
	if err != nil {
		log.Fatal().Err(err).Msg("nats connect failed")
	}
	defer func() { _ = natsConns.NC.Drain() }()
	log.Info().Msg("nats connected")

	// ── Enterprise registry ──────────────────────────────────────────────────
	_ = enterprise.NewRegistry()

	// ── Bifrost / LLM client ─────────────────────────────────────────────────
	log.Info().Msg("initialising bifrost client")
	bf, err := bifrost.New(ctx, cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("bifrost init failed")
	}
	defer bf.Shutdown()

	if bf.Available() {
		log.Info().Str("provider", cfg.DefaultProvider).Str("model", cfg.DefaultModel).Msg("bifrost ready")
	} else {
		log.Warn().Msg("bifrost: no API keys configured — LLM calls will fail at request time")
	}

	// ── Health checker ───────────────────────────────────────────────────────
	checker := health.NewChecker(pool, natsConns.NC)

	// ── Router ───────────────────────────────────────────────────────────────
	r := chi.NewRouter()

	r.Use(cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{
			"Authorization", "Content-Type", "Connect-Protocol-Version",
			"Connect-Timeout-Ms", "Grpc-Timeout", "X-Request-Id",
		},
		ExposedHeaders: []string{"Grpc-Status", "Grpc-Message", "Connect-Status"},
	}).Handler)

	// Infrastructure endpoints
	r.Get("/health", checker.HealthHandler)
	r.Get("/ready", checker.ReadyHandler)

	// ConnectRPC handlers
	connectMux := http.NewServeMux()
	connectMux.Handle(aiv1connect.NewAITaskServiceHandler(
		handler.NewAITaskHandler(bf, log.Logger),
	))

	// Route all ConnectRPC traffic
	r.HandleFunc("/ai.v1.*", connectMux.ServeHTTP)
	r.HandleFunc("/ai.v1.AITaskService/*", connectMux.ServeHTTP)

	// ── HTTP server ──────────────────────────────────────────────────────────
	addr := fmt.Sprintf(":%d", cfg.Port)
	srv := &http.Server{
		Addr:              addr,
		Handler:           h2c.NewHandler(r, &http2.Server{}),
		ReadHeaderTimeout: 10 * time.Second,
		WriteTimeout:      60 * time.Second,
		IdleTimeout:       120 * time.Second,
	}

	go func() {
		log.Info().Str("addr", addr).Msg("llm service listening")
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
	log.Info().Msg("llm service stopped")
}
