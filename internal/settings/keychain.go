package settings

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const keyringService = "com.ignite.app"

type KeychainStore interface {
	Set(service, key, value string) error
	Get(service, key string) (string, error)
	Delete(service, key string) error
}

type goKeychain struct{}

func (g goKeychain) Set(service, key, value string) error {
	return keyring.Set(service, key, value)
}

func (g goKeychain) Get(service, key string) (string, error) {
	return keyring.Get(service, key)
}

func (g goKeychain) Delete(service, key string) error {
	return keyring.Delete(service, key)
}

var DefaultKeychain KeychainStore = goKeychain{}

func SetAPIKey(provider string, key string) error {
	return DefaultKeychain.Set(keyringService, provider, key)
}

func GetAPIKey(provider string) (string, error) {
	secret, err := DefaultKeychain.Get(keyringService, provider)
	if err != nil {
		return "", fmt.Errorf("keychain: no key for %s: %w", provider, err)
	}
	return secret, nil
}

func DeleteAPIKey(provider string) error {
	return DefaultKeychain.Delete(keyringService, provider)
}
