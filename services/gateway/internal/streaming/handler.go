package streaming

import (
	"context"

	"connectrpc.com/connect"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/ApeironFoundation/axle/gateway/internal/hub"

	gatewayv1 "github.com/ApeironFoundation/axle/contracts/generated/go/gateway/v1"
	"github.com/ApeironFoundation/axle/contracts/generated/go/gateway/v1/gatewayv1connect"
	"github.com/google/uuid"
)

// Compile-time interface check.
var _ gatewayv1connect.StreamingServiceHandler = (*Handler)(nil)

// Handler implements the gateway.v1.StreamingService ConnectRPC handler.
type Handler struct {
	hub *hub.Hub
}

// NewHandler returns a StreamingService handler backed by the given Hub.
func NewHandler(h *hub.Hub) *Handler {
	return &Handler{hub: h}
}

// Subscribe implements the server-streaming RPC.
// It subscribes the caller to the hub and streams events until the client
// disconnects or the server shuts down.
func (h *Handler) Subscribe(
	ctx context.Context,
	req *gatewayv1.SubscribeRequest,
	stream *connect.ServerStream[gatewayv1.Event],
) error {
	id := uuid.New().String()
	projectIDs := req.GetProjectIds() // repeated string
	ch, unsub := h.hub.Subscribe(id)
	defer unsub()

	log.Ctx(ctx).Info().
		Str("subscriber_id", id).
		Strs("project_ids", projectIDs).
		Msg("streaming: client connected")

	for {
		select {
		case <-ctx.Done():
			log.Ctx(ctx).Info().Str("subscriber_id", id).Msg("streaming: client disconnected")
			return nil
		case payload, ok := <-ch:
			if !ok {
				return nil
			}
			// Attempt to unmarshal raw NATS payload as a protobuf Event message.
			// NATS publishers must serialize events with proto.Marshal or protojson.
			var event gatewayv1.Event
			if err := proto.Unmarshal(payload, &event); err != nil {
				// Fallback: try JSON (handy for ad-hoc pub from curl / tests).
				if jsonErr := protojson.Unmarshal(payload, &event); jsonErr != nil {
					log.Ctx(ctx).Warn().Err(err).Msg("streaming: dropping unparseable event")
					continue
				}
			}
			if err := stream.Send(&event); err != nil {
				log.Ctx(ctx).Error().Err(err).Msg("streaming: send failed")
				return err
			}
		}
	}
}
