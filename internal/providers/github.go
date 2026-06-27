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

type GitHubCopilotProvider struct {
	token    string
	endpoint string
	client   *http.Client
}

func NewGitHubCopilotProvider(token string) *GitHubCopilotProvider {
	return &GitHubCopilotProvider{
		token:    token,
		endpoint: "https://api.github.com/copilot",
		client:   &http.Client{},
	}
}

type ghCatalogModel struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Publisher   string `json:"publisher"`
	Summary     string `json:"summary"`
	Registry    string `json:"registry"`
	RateLimitTier string `json:"rate_limit_tier"`
	Limits      struct {
		MaxInputTokens  int `json:"max_input_tokens"`
		MaxOutputTokens int `json:"max_output_tokens"`
	} `json:"limits"`
	Capabilities     []string `json:"capabilities"`
	HTMLURL          string   `json:"html_url"`
	SupportedInput  []string `json:"supported_input_modalities"`
	SupportedOutput []string `json:"supported_output_modalities"`
	Tags            []string `json:"tags"`
}

type ghChatRequest struct {
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ghChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type ghStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func (p *GitHubCopilotProvider) ListModels(ctx context.Context) ([]Model, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", "https://models.github.ai/catalog/models", nil)
	if err != nil {
		return nil, fmt.Errorf("github-copilot: models request: %w", err)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github-copilot: models request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return p.defaultModels(), nil
	}

	var catalog []ghCatalogModel
	if err := json.NewDecoder(resp.Body).Decode(&catalog); err != nil {
		return p.defaultModels(), nil
	}

	models := make([]Model, 0, len(catalog))
	for _, m := range catalog {
		if m.ID == "" {
			continue
		}
		models = append(models, Model{ID: m.ID, DisplayName: m.Name})
	}
	if len(models) == 0 {
		return p.defaultModels(), nil
	}
	return models, nil
}

func (p *GitHubCopilotProvider) defaultModels() []Model {
	return []Model{
		{ID: "openai/gpt-4.1", DisplayName: "OpenAI GPT-4.1"},
		{ID: "openai/gpt-4o", DisplayName: "OpenAI GPT-4o"},
		{ID: "openai/gpt-4o-mini", DisplayName: "OpenAI GPT-4o mini"},
		{ID: "openai/gpt-4.1-mini", DisplayName: "OpenAI GPT-4.1 mini"},
		{ID: "openai/gpt-4.1-nano", DisplayName: "OpenAI GPT-4.1 nano"},
	}
}

func (p *GitHubCopilotProvider) Chat(ctx context.Context, model string, messages []Message) (*ChatResponse, error) {
	body := ghChatRequest{Messages: messages, Stream: false}
	reqBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return nil, fmt.Errorf("github-copilot: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("github-copilot: request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("github-copilot: status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var result ghChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("github-copilot: decode: %w", err)
	}
	content := ""
	if len(result.Choices) > 0 {
		content = result.Choices[0].Message.Content
	}
	return &ChatResponse{Content: content, Model: "github-copilot"}, nil
}

func (p *GitHubCopilotProvider) ChatStream(ctx context.Context, model string, messages []Message, onChunk func(string) error) error {
	body := ghChatRequest{Messages: messages, Stream: true}
	reqBody, _ := json.Marshal(body)
	req, err := http.NewRequestWithContext(ctx, "POST", p.endpoint+"/chat/completions", bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("github-copilot: create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+p.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "text/event-stream")

	resp, err := p.client.Do(req)
	if err != nil {
		return fmt.Errorf("github-copilot: request: %w", err)
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
		var chunk ghStreamChunk
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

func (p *GitHubCopilotProvider) ValidateKey(ctx context.Context) error {
	_, err := p.Chat(ctx, "openai/gpt-4.1-mini", []Message{{Role: RoleUser, Content: "ping"}})
	return err
}
