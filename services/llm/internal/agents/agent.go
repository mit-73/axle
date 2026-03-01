// Package agents provides AI agent orchestration logic.
// Agents coordinate multi-step LLM interactions (tool calls, planning, etc.).
package agents

import (
	"context"
	"fmt"

	"github.com/maximhq/bifrost/core/schemas"
	"github.com/rs/zerolog"

	"github.com/ApeironFoundation/axle/llm/internal/bifrostclient"
	"github.com/ApeironFoundation/axle/llm/internal/prompts"
)

// RunRequest describes a single agent run triggered by an AITask.
type RunRequest struct {
	TaskType string
	Payload  []byte
	Model    string
	Provider schemas.ModelProvider
}

// Agent orchestrates a sequence of LLM calls for a given task type.
type Agent struct {
	client *bifrostclient.Client
	log    zerolog.Logger
}

// New returns a new Agent.
func New(client *bifrostclient.Client, log zerolog.Logger) *Agent {
	return &Agent{client: client, log: log}
}

// Run executes the agent loop and streams token chunks to the provided channel.
// The caller is responsible for closing the done channel when finished.
func (a *Agent) Run(ctx context.Context, req RunRequest, out chan<- string) error {
	if !a.client.Available() {
		return fmt.Errorf("bifrost client not available — no API keys configured")
	}

	// Build the system + user messages from the prompt registry.
	messages, err := prompts.BuildMessages(req.TaskType, req.Payload)
	if err != nil {
		return fmt.Errorf("build messages: %w", err)
	}

	model := req.Model
	if model == "" {
		model = a.client.DefaultModel()
	}

	provider := req.Provider
	if provider == "" {
		provider = a.client.DefaultProvider()
	}

	ch, err := a.client.StreamChat(ctx, messages, model, provider)
	if err != nil {
		return fmt.Errorf("stream chat: %w", err)
	}

	for chunk := range ch {
		if chunk == nil || chunk.BifrostChatResponse == nil {
			continue
		}
		for _, choice := range chunk.BifrostChatResponse.Choices {
			if choice.ChatStreamResponseChoice == nil {
				continue
			}
			delta := choice.ChatStreamResponseChoice.Delta
			if delta == nil || delta.Content == nil {
				continue
			}
			out <- *delta.Content
		}
	}

	return nil
}
