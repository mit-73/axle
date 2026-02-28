package handler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	"github.com/google/uuid"
	"github.com/rs/zerolog"

	aiv1 "github.com/ApeironFoundation/axle/contracts/generated/go/ai/v1"
	"github.com/ApeironFoundation/axle/contracts/generated/go/ai/v1/aiv1connect"

	"github.com/ApeironFoundation/axle/llm/internal/agents"
	"github.com/ApeironFoundation/axle/llm/internal/bifrost"
)

// Compile-time interface check.
var _ aiv1connect.AITaskServiceHandler = (*AITaskHandler)(nil)

// AITaskHandler implements ai.v1.AITaskService ConnectRPC methods.
type AITaskHandler struct {
	bifrost *bifrost.Client
	log     zerolog.Logger
}

// NewAITaskHandler creates a new AITaskHandler.
func NewAITaskHandler(bf *bifrost.Client, log zerolog.Logger) *AITaskHandler {
	return &AITaskHandler{bifrost: bf, log: log}
}

// RunAITask streams AI task results back to the caller.
func (h *AITaskHandler) RunAITask(
	ctx context.Context,
	req *aiv1.RunAITaskRequest,
	stream *connect.ServerStream[aiv1.RunAITaskResponse],
) error {
	taskID := uuid.New().String()
	taskType := req.GetType()
	payload := req.GetPayload()

	h.log.Info().
		Str("task_id", taskID).
		Str("type", taskType).
		Str("project_id", req.GetProjectId()).
		Msg("ai task started")

	// Send initial RUNNING status.
	if err := stream.Send(&aiv1.RunAITaskResponse{
		TaskId: taskID,
		Status: aiv1.AITaskStatus_AI_TASK_STATUS_RUNNING,
	}); err != nil {
		return fmt.Errorf("send running status: %w", err)
	}

	// Run agent and stream token chunks.
	agent := agents.New(h.bifrost, h.log)
	tokenCh := make(chan string, 32)
	errCh := make(chan error, 1)

	go func() {
		defer close(tokenCh)
		errCh <- agent.Run(ctx, agents.RunRequest{
			TaskType: taskType,
			Payload:  payload,
		}, tokenCh)
	}()

	for token := range tokenCh {
		if err := stream.Send(&aiv1.RunAITaskResponse{
			TaskId: taskID,
			Status: aiv1.AITaskStatus_AI_TASK_STATUS_RUNNING,
			Chunk:  token,
		}); err != nil {
			return fmt.Errorf("send chunk: %w", err)
		}
	}

	// Check agent error.
	if agentErr := <-errCh; agentErr != nil {
		h.log.Error().Err(agentErr).Str("task_id", taskID).Msg("agent error")
		return stream.Send(&aiv1.RunAITaskResponse{
			TaskId: taskID,
			Status: aiv1.AITaskStatus_AI_TASK_STATUS_FAILED,
			Chunk:  "error: " + agentErr.Error(),
			Done:   true,
		})
	}

	// Send final DONE response.
	return stream.Send(&aiv1.RunAITaskResponse{
		TaskId: taskID,
		Status: aiv1.AITaskStatus_AI_TASK_STATUS_DONE,
		Done:   true,
	})
}
