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
	client   *http.Client
	endpoint string
}

func NewGitHubCopilotProvider(token string) *GitHubCopilotProvider {
	return &GitHubCopilotProvider{
		token: token, endpoint: "https://api.github.com/copilot", client: &http.Client{},
	}
}

type ghCopilotChatRequest struct {
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type ghCopilotChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type ghCopilotStreamChunk struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func (p *GitHubCopilotProvider) Chat(ctx context.Context, model string, messages []Message) (*ChatResponse, error) {
	body := ghCopilotChatRequest{Messages: messages, Stream: false}
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

	var result ghCopilotChatResponse
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
	body := ghCopilotChatRequest{Messages: messages, Stream: true}
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
		var chunk ghCopilotStreamChunk
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

func (p *GitHubCopilotProvider) ListModels(ctx context.Context) ([]Model, error) {
	return []Model{{ID: "github-copilot", DisplayName: "GitHub Copilot"}}, nil
}

func (p *GitHubCopilotProvider) ValidateKey(ctx context.Context) error {
	_, err := p.Chat(ctx, "github-copilot", []Message{{Role: RoleUser, Content: "ping"}})
	return err
}
