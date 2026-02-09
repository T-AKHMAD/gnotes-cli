package cli

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func tokenPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".gnotes", "token"), nil
}

func SaveToken(token string) error {
	p, err := tokenPath()
	if err != nil {
		return err
	}
	dir := filepath.Dir(p)
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return err
	}

	return os.WriteFile(p, []byte(token+"\n"), 0o600)
}

func LoadToken() (string, error) {
	p, err := tokenPath()
	if err != nil {
		return "", err
	}

	b, err := os.ReadFile(p)
	if err != nil {
		return "", err
	}

	token := strings.TrimSpace(string(b))

	if token == "" {
		return "", fmt.Errorf("empty token file: %s", p)
	}

	return token, nil

}
