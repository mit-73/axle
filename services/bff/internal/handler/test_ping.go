package handler

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"

	gatewayv1 "github.com/ApeironFoundation/axle/contracts/go/gateway/v1"
	testv1 "github.com/ApeironFoundation/axle/contracts/go/test/v1"
	"github.com/ApeironFoundation/axle/contracts/go/test/v1/gen_test_v1connect"
)

const (
	testEventsSubject = "axle.events.test.ping"
	testRPCSubject    = "axle.test.ping.rpc"
)

// Compile-time interface check.
var _ gen_test_v1connect.TestServiceHandler = (*TestPingHandler)(nil)

// TestPingHandler is a development-only handler for NATS smoke tests.
type TestPingHandler struct {
	NATS *nats.Conn
}

func (h *TestPingHandler) Ping(
	ctx context.Context,
	req *testv1.PingRequest,
) (*testv1.PingReply, error) {
	if h.NATS == nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("nats is not configured"))
	}

	requestID := req.GetRequestId()
	if requestID == "" {
		requestID = uuid.New().String()
	}

	message := req.GetMessage()
	if message == "" {
		message = "ping"
	}

	ping := &testv1.PingRequest{
		RequestId: requestID,
		Message:   message,
		SentAt:    timestamppb.Now(),
	}

	pingBytes, err := proto.Marshal(ping)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	event := &gatewayv1.Event{
		Id:         requestID,
		Type:       gatewayv1.EventType_EVENT_TYPE_UNSPECIFIED,
		ProjectId:  "dev-only",
		Payload:    pingBytes,
		OccurredAt: timestamppb.Now(),
	}
	eventBytes, err := proto.Marshal(event)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := h.NATS.Publish(testEventsSubject, eventBytes); err != nil {
		return nil, connect.NewError(connect.CodeUnavailable, err)
	}
	log.Info().
		Str("subject", testEventsSubject).
		Str("request_id", requestID).
		Msg("dev-only ping event published")

	msg, err := h.NATS.RequestWithContext(ctx, testRPCSubject, pingBytes)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, nats.ErrTimeout) {
			return nil, connect.NewError(connect.CodeDeadlineExceeded, err)
		}
		return nil, connect.NewError(connect.CodeUnavailable, err)
	}

	var reply testv1.PingReply
	if err := proto.Unmarshal(msg.Data, &reply); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if reply.ReceivedAt == nil {
		reply.ReceivedAt = timestamppb.New(time.Now())
	}
	if reply.RequestId == "" {
		reply.RequestId = requestID
	}

	log.Info().
		Str("subject", testRPCSubject).
		Str("request_id", reply.GetRequestId()).
		Str("responder", reply.GetResponder()).
		Msg("dev-only ping reply received")

	return &reply, nil
}
