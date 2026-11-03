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

// Read is the legacy name-based reader; new sites use [Handle.Read].
// The env name is HGMX_ plus name uppercased with hyphens changed to _.
func Read(name string) (string, error) {
	if !filepath.IsLocal(name) || filepath.Base(name) != name || name == "." {
		return "", fmt.Errorf("credential %q: invalid name", name)
	}
	return read(name, "HGMX_"+strings.ToUpper(strings.ReplaceAll(name, "-", "_")))
}

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
