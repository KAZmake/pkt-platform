package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/KAZmake/pkt-platform/apps/api/assistant-svc/internal/prompt"
	"github.com/KAZmake/pkt-platform/apps/api/assistant-svc/pkg/response"
	anthropic "github.com/anthropics/anthropic-sdk-go"
)

// ChatMessage is a single turn in the conversation history.
type ChatMessage struct {
	Role    string `json:"role"` // "user" | "assistant"
	Content string `json:"content"`
}

// ChatRequest is the body of POST /api/v1/chat.
type ChatRequest struct {
	Message   string        `json:"message"`
	SessionID string        `json:"session_id,omitempty"`
	History   []ChatMessage `json:"history,omitempty"`
}

// ChatResponse is the response body.
type ChatResponse struct {
	Reply     string `json:"reply"`
	SessionID string `json:"session_id,omitempty"`
}

// ChatHandler handles POST /api/v1/chat.
type ChatHandler struct {
	client    *anthropic.Client
	prompt    *prompt.Builder
	model     anthropic.Model
	maxTokens int64
}

func NewChatHandler(client *anthropic.Client, pb *prompt.Builder, model string, maxTokens int) *ChatHandler {
	return &ChatHandler{
		client:    client,
		prompt:    pb,
		model:     anthropic.Model(model),
		maxTokens: int64(maxTokens),
	}
}

func (h *ChatHandler) Chat(w http.ResponseWriter, r *http.Request) {
	var req ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.BadRequest(w, "invalid JSON body")
		return
	}
	if req.Message == "" {
		response.BadRequest(w, "message is required")
		return
	}

	systemPrompt := h.prompt.Get(r.Context())

	// Build the messages list from history + current message.
	msgs := make([]anthropic.MessageParam, 0, len(req.History)+1)
	for _, h := range req.History {
		switch h.Role {
		case "user":
			msgs = append(msgs, anthropic.NewUserMessage(anthropic.NewTextBlock(h.Content)))
		case "assistant":
			msgs = append(msgs, anthropic.NewAssistantMessage(anthropic.NewTextBlock(h.Content)))
		}
	}
	msgs = append(msgs, anthropic.NewUserMessage(anthropic.NewTextBlock(req.Message)))

	result, err := h.client.Messages.New(r.Context(), anthropic.MessageNewParams{
		Model:     anthropic.F(h.model),
		MaxTokens: anthropic.F(h.maxTokens),
		System:    anthropic.F([]anthropic.TextBlockParam{anthropic.NewTextBlock(systemPrompt)}),
		Messages:  anthropic.F(msgs),
	})
	if err != nil {
		slog.Error("anthropic API error", "error", err)
		response.ServiceUnavailable(w, "assistant temporarily unavailable")
		return
	}

	reply := ""
	for _, block := range result.Content {
		if block.Type == "text" {
			reply += block.Text
		}
	}

	response.OK(w, ChatResponse{
		Reply:     reply,
		SessionID: req.SessionID,
	})
}
