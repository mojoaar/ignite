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

type DeepSeekProvider struct {
	apiKey   string
	client   *http.Client
	endpoint string
}

func NewDeepSeekProvider(apiKey string) *DeepSeekProvider {
	return &DeepSeekProvider{
		apiKey: apiKey, endpoint: "https://api.deepseek.com/v1", client: &http.Client{},
	}
}

type deepSeekChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type deepSeekChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type deepSeekStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func (p *DeepSeekProvider) Chat(ctx context.Context, model string, messages []Message) (*ChatResponse, error) {
	body := deepSeekChatRequest{Model: model, Messages: messages, Stream: false}
	reqBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("deepseek: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("deepseek: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("deepseek: status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result deepSeekChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("deepseek: decode: %w", err)
	}
	content := ""
	if len(result.Choices) > 0 {
		content = result.Choices[0].Message.Content
	}
	return &ChatResponse{Content: content, Model: model}, nil
}

func (p *DeepSeekProvider) ChatStream(ctx context.Context, model string, messages []Message, onChunk func(string) error) error {
	body := deepSeekChatRequest{Model: model, Messages: messages, Stream: true}
	reqBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("deepseek: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("deepseek: request: %w", err)
	}
	defer resp.Body.Close()

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
		var chunk deepSeekStreamChunk
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

func (p *DeepSeekProvider) ListModels(ctx context.Context) ([]Model, error) {
	return []Model{
		{ID: "deepseek-v4-flash", DisplayName: "DeepSeek Flash"},
		{ID: "deepseek-v4-pro", DisplayName: "DeepSeek Pro"},
	}, nil
}

func (p *DeepSeekProvider) ValidateKey(ctx context.Context) error {
	_, err := p.Chat(ctx, "deepseek-chat", []Message{{Role: RoleUser, Content: "ping"}})
	return err
}
