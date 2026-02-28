package hub

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
)

// Subscriber is a channel that receives serialised Event bytes.
type Subscriber chan []byte

// Hub manages subscriptions and fan-out of events to connected clients.
// Events arrive from NATS and are broadcast to matching subscribers.
type Hub struct {
	mu   sync.RWMutex
	subs map[string]Subscriber // key: subscriber ID
}

// New creates an empty Hub.
func New() *Hub {
	return &Hub{subs: make(map[string]Subscriber)}
}

// Subscribe registers a new subscriber and returns its channel and an unsubscribe func.
func (h *Hub) Subscribe(id string) (Subscriber, func()) {
	ch := make(Subscriber, 64)
	h.mu.Lock()
	h.subs[id] = ch
	h.mu.Unlock()

	unsub := func() {
		h.mu.Lock()
		delete(h.subs, id)
		close(ch)
		h.mu.Unlock()
	}
	return ch, unsub
}

// Publish fans out payload to all current subscribers.
func (h *Hub) Publish(_ context.Context, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	for id, ch := range h.subs {
		select {
		case ch <- payload:
		default:
			log.Warn().Str("subscriber_id", id).Msg("hub: subscriber buffer full, dropping event")
		}
	}
}
