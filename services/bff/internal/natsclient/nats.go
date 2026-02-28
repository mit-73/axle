package natsclient

import (
	"context"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

// Connections groups the NATS core connection and JetStream context.
type Connections struct {
	NC *nats.Conn
	JS jetstream.JetStream
}

// Connect establishes a NATS connection and initialises a JetStream context.
func Connect(ctx context.Context, natsURL string) (*Connections, error) {
	nc, err := nats.Connect(natsURL,
		nats.Name("axle-bff"),
		nats.Timeout(5*time.Second),
		nats.MaxReconnects(10),
		nats.ReconnectWait(2*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("nats jetstream: %w", err)
	}

	return &Connections{NC: nc, JS: js}, nil
}

// HealthCheck verifies the NATS connection is still alive.
func HealthCheck(_ context.Context, nc *nats.Conn) error {
	if !nc.IsConnected() {
		return fmt.Errorf("nats: not connected (status=%s)", nc.Status())
	}
	return nil
}
