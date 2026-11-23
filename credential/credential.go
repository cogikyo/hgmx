// Package credential declares site secrets and reads their values.
package credential

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrMissing means neither the credential file nor its env fallback was found.
var ErrMissing = errors.New("credential missing")

func read(name, env string) (string, error) {
	if dir := os.Getenv("CREDENTIALS_DIRECTORY"); dir != "" {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err == nil {
			return strings.TrimRight(string(b), "\r\n"), nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("credential %q: %w", name, err)
		}
	}
	if v := os.Getenv(env); v != "" {
		return strings.TrimRight(v, "\r\n"), nil
	}
	return "", fmt.Errorf("%w: %q (or env %s)", ErrMissing, name, env)
}
