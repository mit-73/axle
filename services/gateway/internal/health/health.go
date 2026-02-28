package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/redis/go-redis/v9"
)

// Checker performs health checks on all external dependencies.
type Checker struct {
	redis *redis.Client
	nc    *nats.Conn
}

// NewChecker returns a new Checker.
func NewChecker(redis *redis.Client, nc *nats.Conn) *Checker {
	return &Checker{redis: redis, nc: nc}
}

type status struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}

// HealthHandler returns 200 OK with dependency statuses.
func (c *Checker) HealthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	checks := map[string]string{}
	overall := "ok"

	if err := c.redis.Ping(ctx).Err(); err != nil {
		checks["redis"] = "error: " + err.Error()
		overall = "degraded"
	} else {
		checks["redis"] = "ok"
	}

	if !c.nc.IsConnected() {
		checks["nats"] = "error: not connected"
		overall = "degraded"
	} else {
		checks["nats"] = "ok"
	}

	code := http.StatusOK
	if overall == "degraded" {
		code = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(status{Status: overall, Checks: checks})
}

// ReadyHandler returns 200 OK once all deps are healthy.
func (c *Checker) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	c.HealthHandler(w, r)
}
