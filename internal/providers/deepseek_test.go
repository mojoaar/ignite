package providers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDeepSeek_Chat(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(deepSeekChatResponse{
			Choices: []struct {
				Message struct{ Content string `json:"content"` } `json:"message"`
			}{
				{Message: struct{ Content string `json:"content"` }{Content: "Hello from DeepSeek!"}},
			},
		})
	}))
	defer server.Close()
	p := &DeepSeekProvider{apiKey: "test", endpoint: server.URL, client: server.Client()}
	resp, err := p.Chat(context.Background(), "deepseek-chat", []Message{{Role: RoleUser, Content: "Hi"}})
	if err != nil {
		t.Fatalf("Chat: %v", err)
	}
	if resp.Content != "Hello from DeepSeek!" {
		t.Errorf("unexpected: %s", resp.Content)
	}
}

func TestDeepSeek_ChatStream(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		chunks := []string{
			`data: {"choices":[{"delta":{"content":"Deep"}}]}`,
			`data: {"choices":[{"delta":{"content":"Seek"}}]}`,
			`data: [DONE]`,
		}
		for _, c := range chunks {
			w.Write([]byte(c + "\n"))
		}
	}))
	defer server.Close()
	p := &DeepSeekProvider{apiKey: "test", endpoint: server.URL, client: server.Client()}
	var full string
	err := p.ChatStream(context.Background(), "deepseek-chat", []Message{{Role: RoleUser, Content: "Hi"}}, func(chunk string) error {
		full += chunk
		return nil
	})
	if err != nil {
		t.Fatalf("ChatStream: %v", err)
	}
	if full != "DeepSeek" {
		t.Errorf("unexpected: %q", full)
	}
}

func TestDeepSeek_ChatError401(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer server.Close()
	p := &DeepSeekProvider{apiKey: "test", endpoint: server.URL, client: server.Client()}
	_, err := p.Chat(context.Background(), "deepseek-chat", []Message{{Role: RoleUser, Content: "Hi"}})
	if err == nil {
		t.Error("expected error on 401")
	}
}

func TestDeepSeek_StreamDONE(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Write([]byte("data: {\"choices\":[{\"delta\":{\"content\":\"A\"}}]}\n"))
		w.Write([]byte("data: [DONE]\n"))
	}))
	defer server.Close()
	p := &DeepSeekProvider{apiKey: "test", endpoint: server.URL, client: server.Client()}
	var full string
	err := p.ChatStream(context.Background(), "deepseek-chat", []Message{{Role: RoleUser, Content: "Hi"}}, func(chunk string) error {
		full += chunk
		return nil
	})
	if err != nil {
		t.Fatalf("ChatStream: %v", err)
	}
	if full != "A" {
		t.Errorf("expected A, got %q", full)
	}
}
