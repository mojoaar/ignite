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
	"time"
)

type OpenCodeProvider struct {
	apiKey   string
	client   *http.Client
	endpoint string
}

func NewOpenCodeProvider(apiKey, endpoint string) *OpenCodeProvider {
	return &OpenCodeProvider{
		apiKey:   apiKey,
		endpoint: endpoint,
		client:   &http.Client{Timeout: 120 * time.Second},
	}
}

type openCodeChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type openCodeChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type openCodeStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

type openCodeModelsResponse struct {
	Data []struct {
		ID string `json:"id"`
	} `json:"data"`
}

func (p *OpenCodeProvider) Chat(ctx context.Context, model string, messages []Message) (*ChatResponse, error) {
	body := openCodeChatRequest{Model: model, Messages: messages, Stream: false}
	reqBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("opencode: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("opencode: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("opencode: status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result openCodeChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("opencode: decode: %w", err)
	}
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("opencode: no choices")
	}

	return &ChatResponse{Content: result.Choices[0].Message.Content, Model: model}, nil
}

func (p *OpenCodeProvider) ChatStream(ctx context.Context, model string, messages []Message, onChunk func(string) error) error {
	body := openCodeChatRequest{Model: model, Messages: messages, Stream: true}
	reqBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("opencode: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("opencode: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("opencode: status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "data: ") {
			continue
		}
		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			break
		}
		var chunk openCodeStreamChunk
		if err := json.Unmarshal([]byte(data), &chunk); err != nil {
			continue
		}
		for _, c := range chunk.Choices {
			if c.Delta.Content != "" {
				if err := onChunk(c.Delta.Content); err != nil {
					return err
				}
			}
		}
	}
	return scanner.Err()
}

func (p *OpenCodeProvider) ListModels(ctx context.Context) ([]Model, error) {
	models, err := p.listModelsFromAPI(ctx)
	if err != nil || len(models) == 0 {
		return p.defaultModels(), nil
	}
	return models, nil
}

func (p *OpenCodeProvider) listModelsFromAPI(ctx context.Context) ([]Model, error) {
	req, _ := http.NewRequestWithContext(ctx, "GET", p.endpoint+"/models", nil)
	req.Header.Set("Authorization", "Bearer "+p.apiKey)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("opencode: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("opencode: status %d", resp.StatusCode)
	}

	var result openCodeModelsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("opencode: decode: %w", err)
	}
	out := make([]Model, len(result.Data))
	for i, m := range result.Data {
		out[i] = Model{ID: m.ID, DisplayName: m.ID}
	}
	return out, nil
}

func (p *OpenCodeProvider) defaultModels() []Model {
	if strings.Contains(p.endpoint, "/go/") {
		return p.defaultModelsGo()
	}
	return p.defaultModelsZen()
}

func (p *OpenCodeProvider) defaultModelsGo() []Model {
	return []Model{
		{ID: "glm-5.2", DisplayName: "GLM-5.2"},
		{ID: "glm-5.1", DisplayName: "GLM-5.1"},
		{ID: "kimi-k2.7-code", DisplayName: "Kimi K2.7 Code"},
		{ID: "kimi-k2.6", DisplayName: "Kimi K2.6"},
		{ID: "mimo-v2.5", DisplayName: "MiMo-V2.5"},
		{ID: "mimo-v2.5-pro", DisplayName: "MiMo-V2.5-Pro"},
		{ID: "minimax-m3", DisplayName: "MiniMax M3"},
		{ID: "minimax-m2.7", DisplayName: "MiniMax M2.7"},
		{ID: "qwen3.7-max", DisplayName: "Qwen3.7 Max"},
		{ID: "qwen3.7-plus", DisplayName: "Qwen3.7 Plus"},
		{ID: "qwen3.6-plus", DisplayName: "Qwen3.6 Plus"},
		{ID: "deepseek-v4-pro", DisplayName: "DeepSeek V4 Pro"},
		{ID: "deepseek-v4-flash", DisplayName: "DeepSeek V4 Flash"},
	}
}

func (p *OpenCodeProvider) defaultModelsZen() []Model {
	return []Model{
		{ID: "gpt-5.5", DisplayName: "GPT-5.5"},
		{ID: "gpt-5.5-pro", DisplayName: "GPT-5.5 Pro"},
		{ID: "gpt-5.4", DisplayName: "GPT-5.4"},
		{ID: "gpt-5.4-pro", DisplayName: "GPT-5.4 Pro"},
		{ID: "gpt-5.4-mini", DisplayName: "GPT-5.4 Mini"},
		{ID: "gpt-5.4-nano", DisplayName: "GPT-5.4 Nano"},
		{ID: "gpt-5.3-codex", DisplayName: "GPT-5.3 Codex"},
		{ID: "gpt-5.3-codex-spark", DisplayName: "GPT-5.3 Codex Spark"},
		{ID: "gpt-5.2", DisplayName: "GPT-5.2"},
		{ID: "gpt-5.1", DisplayName: "GPT-5.1"},
		{ID: "gpt-5.1-codex", DisplayName: "GPT-5.1 Codex"},
		{ID: "gpt-5", DisplayName: "GPT-5"},
		{ID: "gpt-5-codex", DisplayName: "GPT-5 Codex"},
		{ID: "gpt-5-nano", DisplayName: "GPT-5 Nano"},
		{ID: "claude-fable-5", DisplayName: "Claude Fable 5"},
		{ID: "claude-opus-4-8", DisplayName: "Claude Opus 4.8"},
		{ID: "claude-opus-4-5", DisplayName: "Claude Opus 4.5"},
		{ID: "claude-sonnet-4-5", DisplayName: "Claude Sonnet 4.5"},
		{ID: "claude-sonnet-4", DisplayName: "Claude Sonnet 4"},
		{ID: "claude-haiku-4-5", DisplayName: "Claude Haiku 4.5"},
		{ID: "gemini-3.5-flash", DisplayName: "Gemini 3.5 Flash"},
		{ID: "gemini-3.1-pro", DisplayName: "Gemini 3.1 Pro"},
		{ID: "gemini-3-flash", DisplayName: "Gemini 3 Flash"},
		{ID: "grok-build-0.1", DisplayName: "Grok Build 0.1"},
		{ID: "deepseek-v4-pro", DisplayName: "DeepSeek V4 Pro"},
		{ID: "deepseek-v4-flash", DisplayName: "DeepSeek V4 Flash"},
		{ID: "glm-5.2", DisplayName: "GLM-5.2"},
		{ID: "glm-5.1", DisplayName: "GLM-5.1"},
		{ID: "glm-5", DisplayName: "GLM-5"},
	}
}

func (p *OpenCodeProvider) ValidateKey(ctx context.Context) error {
	_, err := p.ListModels(ctx)
	return err
}
