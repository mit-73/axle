package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	gatewayv1 "github.com/ApeironFoundation/axle/contracts/go/gateway/v1"
	testv1 "github.com/ApeironFoundation/axle/contracts/go/test/v1"
)

const (
	eventsSubject = "axle.events.test.ping"
	rpcSubject    = "axle.test.ping.rpc"
)

// StartDevPingSubscriptions registers development-only NATS subscriptions that
// validate pub/sub and request-reply inter-service communication.
func StartDevPingSubscriptions(_ context.Context, nc *nats.Conn, logger zerolog.Logger) ([]*nats.Subscription, error) {
	eventsSub, err := nc.Subscribe(eventsSubject, func(msg *nats.Msg) {
		var event gatewayv1.Event
		if err := proto.Unmarshal(msg.Data, &event); err != nil {
			logger.Warn().Err(err).Str("subject", msg.Subject).Msg("dev-only ping event: unmarshal gateway event failed")
			return
		}
		var ping testv1.PingRequest
		if err := proto.Unmarshal(event.GetPayload(), &ping); err != nil {
			logger.Warn().Err(err).Str("subject", msg.Subject).Msg("dev-only ping event: unmarshal ping payload failed")
			return
		}

		logger.Info().
			Str("subject", msg.Subject).
			Str("request_id", ping.GetRequestId()).
			Str("message", ping.GetMessage()).
			Msg("dev-only ping event received by llm")
	})
	if err != nil {
		return nil, fmt.Errorf("subscribe %s: %w", eventsSubject, err)
	}

	rpcSub, err := nc.Subscribe(rpcSubject, func(msg *nats.Msg) {
		var ping testv1.PingRequest
		if err := proto.Unmarshal(msg.Data, &ping); err != nil {
			logger.Warn().Err(err).Str("subject", msg.Subject).Msg("dev-only ping rpc: bad request payload")
			return
		}

		reply := &testv1.PingReply{
			RequestId:  ping.GetRequestId(),
			Message:    "pong from llm",
			Responder:  "llm",
			ReceivedAt: timestamppb.Now(),
		}
		payload, err := proto.Marshal(reply)
		if err != nil {
			logger.Error().Err(err).Str("subject", msg.Subject).Msg("dev-only ping rpc: marshal reply failed")
			return
		}

		if err := msg.Respond(payload); err != nil {
			logger.Error().Err(err).Str("subject", msg.Subject).Msg("dev-only ping rpc: respond failed")
			return
		}

		logger.Info().
			Str("subject", msg.Subject).
			Str("request_id", ping.GetRequestId()).
			Msg("dev-only ping rpc request handled by llm")
	})
	if err != nil {
		_ = eventsSub.Unsubscribe()
		return nil, fmt.Errorf("subscribe %s: %w", rpcSubject, err)
	}

	return []*nats.Subscription{eventsSub, rpcSub}, nil
}
