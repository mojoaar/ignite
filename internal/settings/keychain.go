package settings

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const keyringService = "com.ignite.app"

func SetAPIKey(provider string, key string) error {
	return keyring.Set(keyringService, provider, key)
}

func GetAPIKey(provider string) (string, error) {
	secret, err := keyring.Get(keyringService, provider)
	if err != nil {
		return "", fmt.Errorf("keychain: no key for %s: %w", provider, err)
	}
	return secret, nil
}

func DeleteAPIKey(provider string) error {
	return keyring.Delete(keyringService, provider)
}
