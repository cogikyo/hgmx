package auth

import (
	"io"
	"testing"

	"github.com/cogikyo/hgmx/cli"
)

func TestCommandExit(t *testing.T) {
	a := setup(t)
	for _, tc := range []struct {
		args []string
		code int
	}{
		{nil, 2},
		{[]string{"rotate", "x@example.com"}, 2},
		{[]string{"link"}, 2},
		{[]string{"link", "--nope", "x@example.com"}, 2},
		{[]string{"link", "x@example.com"}, 2},
		{[]string{"link", "--origin", "https://other.example.com", "x@example.com"}, 2},
		{[]string{"link", "--origin", "not a url", "x@example.com"}, 2},
		{[]string{"link", "--origin", apex, "nobody"}, 1},
		{[]string{"link", "--origin", apex, "x@example.com"}, 0},
	} {
		err := Command(t.Context(), a.db, config(), tc.args, io.Discard)
		if got := cli.Exit(err); got != tc.code {
			t.Errorf("%q: exit %d (%v), want %d", tc.args, got, err, tc.code)
		}
	}
}
