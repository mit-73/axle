// Package prompts manages system and user prompt templates for each AI task type.
package prompts

import (
	"encoding/json"
	"fmt"

	"github.com/maximhq/bifrost/core/schemas"
)

// systemPrompts maps task type → system prompt text.
var systemPrompts = map[string]string{
	"summarise":   "You are a concise summarisation assistant. Summarise the provided text clearly and briefly.",
	"explain":     "You are a helpful technical tutor. Explain the provided concept clearly.",
	"code_review": "You are a senior software engineer. Review the provided code and give actionable feedback.",
	"translate":   "You are a professional translator. Translate the provided text accurately.",
	// default fallback handled in BuildMessages
}

// payloadInput is the expected shape of the JSON payload field.
type payloadInput struct {
	Text  string `json:"text"`
	Query string `json:"query"`
}

// BuildMessages constructs the messages slice for a given task type and raw payload.
func BuildMessages(taskType string, payload []byte) ([]schemas.ChatMessage, error) {
	var input payloadInput
	if len(payload) > 0 {
		if err := json.Unmarshal(payload, &input); err != nil {
			// Treat the raw bytes as plain text if not valid JSON.
			input.Text = string(payload)
		}
	}

	text := input.Text
	if text == "" {
		text = input.Query
	}
	if text == "" {
		return nil, fmt.Errorf("payload contains no 'text' or 'query' field")
	}

	sysPrompt, ok := systemPrompts[taskType]
	if !ok {
		sysPrompt = "You are a helpful AI assistant."
	}

	sysRole := schemas.ChatMessageRoleSystem
	userRole := schemas.ChatMessageRoleUser

	sysContent := &schemas.ChatMessageContent{ContentStr: &sysPrompt}
	userContent := &schemas.ChatMessageContent{ContentStr: &text}

	return []schemas.ChatMessage{
		{Role: sysRole, Content: sysContent},
		{Role: userRole, Content: userContent},
	}, nil
}
