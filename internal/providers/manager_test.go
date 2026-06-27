package providers

import (
	"testing"
)

func TestManager_RegisterAndGet(t *testing.T) {
	m := NewManager()
	p := NewDeepSeekProvider("test-key")
	m.Register("deepseek", p)

	got, err := m.Get("deepseek")
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got != p {
		t.Error("Get returned wrong provider")
	}
}

func TestManager_GetMissing(t *testing.T) {
	m := NewManager()
	_, err := m.Get("nonexistent")
	if err == nil {
		t.Error("expected error for missing provider")
	}
}
