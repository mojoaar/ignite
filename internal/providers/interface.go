package providers

import "context"

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
	RoleSystem    Role = "system"
)

type Message struct {
	Role    Role   `json:"role"`
	Content string `json:"content"`
}

type Model struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name"`
}

type ChatResponse struct {
	Content string `json:"content"`
	Model   string `json:"model"`
}

type LLMProvider interface {
	Chat(ctx context.Context, model string, messages []Message) (*ChatResponse, error)
	ChatStream(ctx context.Context, model string, messages []Message, onChunk func(chunk string) error) error
	ListModels(ctx context.Context) ([]Model, error)
	ValidateKey(ctx context.Context) error
}
