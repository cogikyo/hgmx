package icons

import (
	"context"
	"io/fs"
	"strings"
	"testing"
)

func TestRender(t *testing.T) {
	ctx := context.Background()
	for _, n := range All() {
		var b strings.Builder
		if err := n.Render(ctx, &b); err != nil {
			t.Fatalf("%s: %v", n, err)
		}
		html := b.String()
		if !strings.HasPrefix(html, "<svg") || !strings.HasSuffix(html, "</svg>") {
			t.Errorf("%s: not one svg element: %.60q", n, html)
		}
		_, err := fs.Stat(assets, "assets/anim/"+string(n)+".js")
		if got, animated := strings.Contains(html, ` data-icon="`+string(n)+`"`), err == nil; got != animated {
			t.Errorf("%s: data-icon present %v, animated %v", n, got, animated)
		}
		if glyphs[n].tags == nil {
			t.Errorf("%s: no tags", n)
		}
	}
	if err := Name("no-such-icon").Render(ctx, &strings.Builder{}); err == nil {
		t.Error("unknown name rendered without error")
	}
}
