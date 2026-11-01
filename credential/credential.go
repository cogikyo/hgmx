// Package credential reads credential files with environment fallbacks.
package credential

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ErrMissing is wrapped by Read when the file and env fallback are missing.
// Test with errors.Is to disable optional features; treat other errors as fatal.
var ErrMissing = errors.New("credential missing")

// Read reads $CREDENTIALS_DIRECTORY/<name>, trimming trailing CR and LF.
// If the directory is unset or the file is missing, it reads an env var.
// The env name is HGMX_ plus name uppercased with hyphens changed to _.
// Other file errors prevent fallback; even an empty file takes priority.
// Names must be single local path components; empty env values are missing.
func Read(name string) (string, error) {
	if !filepath.IsLocal(name) || filepath.Base(name) != name || name == "." {
		return "", fmt.Errorf("credential %q: invalid name", name)
	}
	if dir := os.Getenv("CREDENTIALS_DIRECTORY"); dir != "" {
		b, err := os.ReadFile(filepath.Join(dir, name))
		if err == nil {
			return strings.TrimRight(string(b), "\r\n"), nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("credential %q: %w", name, err)
		}
	}
	env := "HGMX_" + strings.ToUpper(strings.ReplaceAll(name, "-", "_"))
	if v := os.Getenv(env); v != "" {
		return strings.TrimRight(v, "\r\n"), nil
	}
	return "", fmt.Errorf("%w: %q (or env %s)", ErrMissing, name, env)
}
