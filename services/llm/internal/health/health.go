package health

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/nats-io/nats.go"
)

// Checker performs health checks on all external dependencies.
type Checker struct {
	db *pgxpool.Pool
	nc *nats.Conn
}

// NewChecker returns a new Checker.
func NewChecker(db *pgxpool.Pool, nc *nats.Conn) *Checker {
	return &Checker{db: db, nc: nc}
}

type statusResponse struct {
	Status string            `json:"status"`
	Checks map[string]string `json:"checks"`
}

// HealthHandler returns 200 OK with dependency statuses (or 503 if degraded).
func (c *Checker) HealthHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
	defer cancel()

	checks := make(map[string]string, 2)
	overall := "ok"

	if c.db != nil {
		if err := c.db.Ping(ctx); err != nil {
			checks["postgres"] = "error: " + err.Error()
			overall = "degraded"
		} else {
			checks["postgres"] = "ok"
		}
	} else {
		checks["postgres"] = "unconfigured"
	}

	if c.nc != nil && c.nc.IsConnected() {
		checks["nats"] = "ok"
	} else {
		checks["nats"] = "error: not connected"
		overall = "degraded"
	}

	code := http.StatusOK
	if overall == "degraded" {
		code = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(statusResponse{Status: overall, Checks: checks})
}

// ReadyHandler is an alias for HealthHandler — ready once all deps are healthy.
func (c *Checker) ReadyHandler(w http.ResponseWriter, r *http.Request) {
	c.HealthHandler(w, r)
}
