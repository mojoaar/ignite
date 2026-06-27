package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestClaude_Chat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(claudeResponse{Content: []struct{ Text string `json:"text"` }{
			{Text: "Hello from Claude!"},
		}})
	}))
	defer server.Close()
	p := &ClaudeProvider{apiKey: "test", endpoint: server.URL, client: server.Client()}
	resp, err := p.Chat(context.Background(), "claude-sonnet-4", []Message{{Role: RoleUser, Content: "Hi"}})
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}
	if resp.Content != "Hello from Claude!" {
		t.Errorf("unexpected: %s", resp.Content)
	}
}

func TestClaude_ListModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(claudeModelsResponse{Data: []struct {
			ID          string `json:"id"`
			DisplayName string `json:"display_name"`
		}{{ID: "claude-sonnet-4", DisplayName: "Claude Sonnet 4"}}})
	}))
	defer server.Close()
	p := &ClaudeProvider{apiKey: "test", endpoint: server.URL, client: server.Client()}
	models, err := p.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels: %v", err)
	}
	if len(models) != 1 {
		t.Errorf("expected 1, got %d", len(models))
	}
}

func TestClaude_ChatStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		events := []string{
			`data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"Hello"}}`,
			`data: {"type":"content_block_delta","delta":{"type":"text_delta","text":" Claude"}}`,
		}
		for _, e := range events {
			w.Write([]byte(e + "\n"))
		}
	}))
	defer server.Close()
	p := &ClaudeProvider{apiKey: "test", endpoint: server.URL, client: server.Client()}
	var full string
	err := p.ChatStream(context.Background(), "claude-sonnet-4", []Message{{Role: RoleUser, Content: "Hi"}}, func(chunk string) error {
		full += chunk
		return nil
	})
	if err != nil {
		t.Fatalf("ChatStream: %v", err)
	}
	if full != "Hello Claude" {
		t.Errorf("unexpected: %q", full)
	}
}
