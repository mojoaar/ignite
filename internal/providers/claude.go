package providers

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type ClaudeProvider struct {
	apiKey   string
	client   *http.Client
	endpoint string
}

func NewClaudeProvider(apiKey string) *ClaudeProvider {
	return &ClaudeProvider{
		apiKey: apiKey, endpoint: "https://api.anthropic.com/v1/messages", client: &http.Client{},
	}
}

type claudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type claudeRequest struct {
	Model     string          `json:"model"`
	MaxTokens int             `json:"max_tokens"`
	Messages  []claudeMessage `json:"messages"`
	Stream    bool            `json:"stream"`
}

type claudeResponse struct {
	Content []struct {
		Text string `json:"text"`
	} `json:"content"`
}

type claudeStreamEvent struct {
	Type  string `json:"type"`
	Delta struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta"`
}

func convertToClaude(messages []Message) []claudeMessage {
	out := make([]claudeMessage, len(messages))
	for i, m := range messages {
		role := string(m.Role)
		if role == "system" {
			role = "user"
		}
		out[i] = claudeMessage{Role: role, Content: m.Content}
	}
	return out
}

func (p *ClaudeProvider) Chat(ctx context.Context, model string, messages []Message) (*ChatResponse, error) {
	body := claudeRequest{Model: model, MaxTokens: 4096, Messages: convertToClaude(messages), Stream: false}
	reqBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("claude: create request: %w", err)
	}
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("claude: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("claude: status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result claudeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("claude: decode: %w", err)
	}
	content := ""
	if len(result.Content) > 0 {
		content = result.Content[0].Text
	}
	return &ChatResponse{Content: content, Model: model}, nil
}

func (p *ClaudeProvider) ChatStream(ctx context.Context, model string, messages []Message, onChunk func(string) error) error {
	body := claudeRequest{Model: model, MaxTokens: 4096, Messages: convertToClaude(messages), Stream: true}
	reqBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("claude: create request: %w", err)
	}
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("claude: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("claude: status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		var event claudeStreamEvent
		if err := json.Unmarshal([]byte(data), &event); err != nil {
			continue
		}
		if event.Type == "content_block_delta" && event.Delta.Type == "text_delta" {
			if err := onChunk(event.Delta.Text); err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}

type claudeModelsResponse struct {
	Data []struct {
		ID          string `json:"id"`
		DisplayName string `json:"display_name"`
	} `json:"data"`
}

func (p *ClaudeProvider) ListModels(ctx context.Context) ([]Model, error) {
	models, err := p.listModelsFromAPI(ctx)
	if err != nil || len(models) == 0 {
		return p.defaultModels(), nil
	}
	return models, nil
}

func (p *ClaudeProvider) listModelsFromAPI(ctx context.Context) ([]Model, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", p.endpoint+"/models?limit=100", nil)
	if err != nil {
		return nil, fmt.Errorf("claude: create request: %w", err)
	}
	req.Header.Set("x-api-key", p.apiKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("claude: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("claude: status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result claudeModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("claude: decode: %w", err)
	}
	out := make([]Model, len(result.Data))
	for i, m := range result.Data {
		out[i] = Model{ID: m.ID, DisplayName: m.DisplayName}
	}
	return out, nil
}

func (p *ClaudeProvider) defaultModels() []Model {
	return []Model{
		{ID: "claude-sonnet-4-20250514", DisplayName: "Claude Sonnet 4"},
		{ID: "claude-opus-4-20250514", DisplayName: "Claude Opus 4"},
		{ID: "claude-haiku-3.5", DisplayName: "Claude Haiku 3.5"},
	}
}

func (p *ClaudeProvider) ValidateKey(ctx context.Context) error {
	_, err := p.Chat(ctx, "claude-haiku-3.5", []Message{{Role: RoleUser, Content: "ping"}})
	return err
}
