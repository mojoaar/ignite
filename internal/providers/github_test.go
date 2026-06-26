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
		json.NewEncoder(w).Encode(ghCopilotChatResponse{
			Choices: []struct {
				Message struct{ Content string `json:"content"` } `json:"message"`
			}{
				{Message: struct{ Content string `json:"content"` }{Content: "Hello from Copilot!"}},
			},
		})
	}))
	defer server.Close()
	p := &GitHubCopilotProvider{token: "test-token", endpoint: server.URL, client: server.Client()}
	resp, err := p.Chat(context.Background(), "github-copilot", []Message{{Role: RoleUser, Content: "Hi"}})
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}
	if resp.Content != "Hello from Copilot!" {
		t.Errorf("unexpected: %s", resp.Content)
	}
}
