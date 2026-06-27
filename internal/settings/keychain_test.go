package settings

import (
	"testing"
)

type mockKeychain struct {
	data map[string]string
}

func (m *mockKeychain) Set(service, key, value string) error {
	if m.data == nil {
		m.data = make(map[string]string)
	}
	m.data[service+"/"+key] = value
	return nil
}

func (m *mockKeychain) Get(service, key string) (string, error) {
	v, ok := m.data[service+"/"+key]
	if !ok {
		return "", nil
	}
	return v, nil
}

func (m *mockKeychain) Delete(service, key string) error {
	delete(m.data, service+"/"+key)
	return nil
}

func TestKeychain_SetAndGet(t *testing.T) {
	m := &mockKeychain{}
	DefaultKeychain = m
	defer func() { DefaultKeychain = goKeychain{} }()

	if err := SetAPIKey("opencode-go", "sk-test"); err != nil {
		t.Fatalf("SetAPIKey: %v", err)
	}
	val, err := GetAPIKey("opencode-go")
	if err == nil && val == "sk-test" {
		return
	}
	if err != nil {
		t.Fatalf("GetAPIKey: %v", err)
	}
	if val != "sk-test" {
		t.Errorf("expected sk-test, got %s", val)
	}
}

func TestKeychain_Delete(t *testing.T) {
	m := &mockKeychain{}
	DefaultKeychain = m
	defer func() { DefaultKeychain = goKeychain{} }()

	SetAPIKey("deepseek", "sk-ds")
	DeleteAPIKey("deepseek")
	val, _ := GetAPIKey("deepseek")
	if val != "" {
		t.Errorf("expected empty after delete, got %s", val)
	}
}

func TestKeychain_SharedOpenCode(t *testing.T) {
	m := &mockKeychain{}
	DefaultKeychain = m
	defer func() { DefaultKeychain = goKeychain{} }()

	SetAPIKey("opencode-go", "sk-shared")
	val, _ := GetAPIKey("opencode-go")
	if val != "sk-shared" {
		t.Errorf("expected shared key, got %s", val)
	}
}
