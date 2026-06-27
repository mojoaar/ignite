package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGitHubCopilot_Chat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Error("missing auth")
		}
		json.NewEncoder(w).Encode(ghChatResponse{
			Choices: []struct {
				Message struct{ Content string `json:"content"` } `json:"message"`
			}{{Message: struct{ Content string `json:"content"` }{Content: "Hello from Copilot!"}}},
		})
	}))
	defer server.Close()
	p := &GitHubCopilotProvider{token: "test-token", endpoint: server.URL, client: server.Client()}
	resp, err := p.Chat(context.Background(), "openai/gpt-4o", []Message{{Role: RoleUser, Content: "Hi"}})
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}
	if resp.Content != "Hello from Copilot!" {
		t.Errorf("unexpected: %s", resp.Content)
	}
}

func TestGitHubCopilot_ListModels(t *testing.T) {
	p := NewGitHubCopilotProvider("test-token")
	models, err := p.ListModels(context.Background())
	if err != nil {
		t.Fatalf("ListModels: %v", err)
	}
	if len(models) < 2 {
		t.Errorf("expected at least 2 models, got %d", len(models))
	}
}
