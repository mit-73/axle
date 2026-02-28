// Package bifrost wraps the maximhq/bifrost SDK and provides a thin Client
// used by the LLM service handlers.
package bifrost

import (
	"context"
	"fmt"

	bifrost "github.com/maximhq/bifrost/core"
	"github.com/maximhq/bifrost/core/schemas"

	"github.com/ApeironFoundation/axle/llm/internal/config"
)

// ── Account implementation ────────────────────────────────────────────────────

// localAccount implements schemas.Account, reading API keys from Config.
type localAccount struct {
	cfg *config.Config
}

// GetConfiguredProviders returns only providers that have an API key set.
func (a *localAccount) GetConfiguredProviders() ([]schemas.ModelProvider, error) {
	var providers []schemas.ModelProvider
	if a.cfg.OpenAIAPIKey != "" {
		providers = append(providers, schemas.OpenAI)
	}
	if a.cfg.AnthropicAPIKey != "" {
		providers = append(providers, schemas.Anthropic)
	}
	return providers, nil
}

// GetKeysForProvider returns key(s) for the given provider.
func (a *localAccount) GetKeysForProvider(_ context.Context, p schemas.ModelProvider) ([]schemas.Key, error) {
	switch p {
	case schemas.OpenAI:
		if a.cfg.OpenAIAPIKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY not set")
		}
		return []schemas.Key{
			{
				Value:  *schemas.NewEnvVar(a.cfg.OpenAIAPIKey),
				Weight: 1.0,
			},
		}, nil
	case schemas.Anthropic:
		if a.cfg.AnthropicAPIKey == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
		}
		return []schemas.Key{
			{
				Value:  *schemas.NewEnvVar(a.cfg.AnthropicAPIKey),
				Weight: 1.0,
			},
		}, nil
	default:
		return nil, fmt.Errorf("provider %q not configured", p)
	}
}

// GetConfigForProvider returns network/concurrency config for a provider.
func (a *localAccount) GetConfigForProvider(p schemas.ModelProvider) (*schemas.ProviderConfig, error) {
	switch p {
	case schemas.OpenAI, schemas.Anthropic:
		return &schemas.ProviderConfig{
			ConcurrencyAndBufferSize: schemas.ConcurrencyAndBufferSize{
				Concurrency: 10,
				BufferSize:  100,
			},
		}, nil
	default:
		return nil, fmt.Errorf("provider %q not configured", p)
	}
}

// ── Client ────────────────────────────────────────────────────────────────────

// Client wraps a bifrost.Bifrost instance.
// When no API keys are configured, Client is created in a "disabled" state
// (bf == nil) and Available() returns false.
type Client struct {
	bf  *bifrost.Bifrost
	cfg *config.Config
}

// New initialises a bifrost Client from configuration.
// If no providers are configured (no API keys), returns a disabled client
// instead of an error so the service can start without LLM credentials.
func New(ctx context.Context, cfg *config.Config) (*Client, error) {
	acct := &localAccount{cfg: cfg}

	providers, err := acct.GetConfiguredProviders()
	if err != nil {
		return nil, fmt.Errorf("bifrost: get providers: %w", err)
	}

	if len(providers) == 0 {
		// Graceful degradation — no keys set, LLM calls will return an error at
		// request time rather than at startup.
		return &Client{bf: nil, cfg: cfg}, nil
	}

	bf, err := bifrost.Init(ctx, schemas.BifrostConfig{
		Account: acct,
		Logger:  bifrost.NewDefaultLogger(schemas.LogLevelInfo),
	})
	if err != nil {
		return nil, fmt.Errorf("bifrost init: %w", err)
	}

	return &Client{bf: bf, cfg: cfg}, nil
}

// Available reports whether the bifrost client is active (i.e. at least one
// API key was configured at startup).
func (c *Client) Available() bool {
	return c.bf != nil
}

// Shutdown gracefully shuts down the underlying bifrost instance.
func (c *Client) Shutdown() {
	if c.bf != nil {
		c.bf.Shutdown()
	}
}

// StreamChat sends a chat completion request and returns a channel of stream
// chunks. The caller must drain the channel to completion.
// Returns an error if no providers are configured.
func (c *Client) StreamChat(
	ctx context.Context,
	messages []schemas.ChatMessage,
	model string,
	provider schemas.ModelProvider,
) (chan *schemas.BifrostStreamChunk, error) {
	if !c.Available() {
		return nil, fmt.Errorf("bifrost: no providers configured — set OPENAI_API_KEY or ANTHROPIC_API_KEY")
	}

	bCtx, cancel := schemas.NewBifrostContextWithCancel(ctx)
	_ = cancel // channel close signals completion — keep cancel for GC

	req := &schemas.BifrostChatRequest{
		Provider: provider,
		Model:    model,
		Input:    messages,
	}

	ch, err := c.bf.ChatCompletionStreamRequest(bCtx, req)
	if err != nil {
		return nil, fmt.Errorf("bifrost stream: %s", bifrost.GetErrorMessage(err))
	}

	return ch, nil
}

// DefaultProvider returns the configured default ModelProvider constant.
func (c *Client) DefaultProvider() schemas.ModelProvider {
	switch c.cfg.DefaultProvider {
	case "anthropic":
		return schemas.Anthropic
	case "gemini":
		return schemas.Gemini
	case "mistral":
		return schemas.Mistral
	default:
		return schemas.OpenAI
	}
}

// DefaultModel returns the configured default model name.
func (c *Client) DefaultModel() string {
	return c.cfg.DefaultModel
}
