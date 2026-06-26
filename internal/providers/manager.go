package providers

import "fmt"

type Manager struct {
	providers map[string]LLMProvider
}

func NewManager() *Manager {
	return &Manager{providers: make(map[string]LLMProvider)}
}

func (m *Manager) Register(name string, p LLMProvider) {
	m.providers[name] = p
}

func (m *Manager) Get(name string) (LLMProvider, error) {
	p, ok := m.providers[name]
	if !ok {
		return nil, fmt.Errorf("provider %s not configured", name)
	}
	return p, nil
}
