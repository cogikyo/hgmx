package cli

import (
	"errors"
	"fmt"
	"testing"
)

func TestExit(t *testing.T) {
	usage := Usage("usage: migrate up")
	for _, tc := range []struct {
		err  error
		code int
	}{
		{nil, 0},
		{errors.New("PGDATABASE is unset"), 1},
		{usage, 2},
		{fmt.Errorf("migrate: %w", usage), 2},
		{Usage("site flags: %w", errors.New("bad flag")), 2},
	} {
		if got := Exit(tc.err); got != tc.code {
			t.Errorf("Exit(%v) = %d, want %d", tc.err, got, tc.code)
		}
	}
	if usage.Error() != "usage: migrate up" {
		t.Errorf("Usage text %q, want the formatted message alone", usage.Error())
	}
}
