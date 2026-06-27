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

type OpenCodeProvider struct {
	apiKey   string
	client   *http.Client
	endpoint string
}

func NewOpenCodeProvider(apiKey, endpoint string) *OpenCodeProvider {
	return &OpenCodeProvider{
		apiKey: apiKey, endpoint: endpoint, client: &http.Client{},
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
	models := make([]Model, len(result.Data))
	for i, m := range result.Data {
		models[i] = Model{ID: m.ID, DisplayName: m.ID}
	}
	return models, nil
}

func (p *OpenCodeProvider) ValidateKey(ctx context.Context) error {
	_, err := p.ListModels(ctx)
	return err
}
