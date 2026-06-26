package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenCode_Chat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(openCodeChatResponse{
			Choices: []struct {
				Message struct{ Content string `json:"content"` } `json:"message"`
			}{
				{Message: struct{ Content string `json:"content"` }{Content: "Hello from OpenCode!"}},
			},
		})
	}))
	defer server.Close()

	p := &OpenCodeProvider{apiKey: "test", endpoint: server.URL, client: server.Client()}
	resp, err := p.Chat(context.Background(), "gpt-4o", []Message{{Role: RoleUser, Content: "Hi"}})
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}
	if resp.Content != "Hello from OpenCode!" {
		t.Errorf("unexpected: %s", resp.Content)
	}
}

func TestOpenCode_ChatStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		chunks := []string{
			`data: {"choices":[{"delta":{"content":"Hello"}}]}`,
			`data: {"choices":[{"delta":{"content":" world"}}]}`,
			`data: [DONE]`,
		}
		for _, c := range chunks {
			w.Write([]byte(c + "\n"))
		}
	}))
	defer server.Close()

	p := &OpenCodeProvider{apiKey: "test", endpoint: server.URL, client: server.Client()}
	var full string
	err := p.ChatStream(context.Background(), "gpt-4o", []Message{{Role: RoleUser, Content: "Hi"}}, func(chunk string) error {
		full += chunk
		return nil
	})
	if err != nil {
		t.Fatalf("ChatStream: %v", err)
	}
	if full != "Hello world" {
		t.Errorf("unexpected: %q", full)
	}
}

func TestOpenCode_ListModels(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(openCodeModelsResponse{Data: []struct{ ID string `json:"id"` }{
			{ID: "gpt-4o"}, {ID: "gpt-4o-mini"},
		}})
	}))
	defer server.Close()

	p := &OpenCodeProvider{apiKey: "test", endpoint: server.URL, client: server.Client()}
	models, err := p.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels: %v", err)
	}
	if len(models) != 2 {
		t.Errorf("expected 2, got %d", len(models))
	}
}
