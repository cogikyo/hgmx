package icons

import (
	"context"
	"io/fs"
	"slices"
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
		_, err := fs.Stat(FS, "anim/"+string(n)+".js")
		if got, animated := strings.Contains(html, ` data-icon="`+string(n)+`"`), err == nil; got != animated {
			t.Errorf("%s: data-icon present %v, animated %v", n, got, animated)
		}
		if n.Tags() == nil {
			t.Errorf("%s: no tags", n)
		}
	}
	if err := Name("no-such-icon").Render(ctx, &strings.Builder{}); err == nil {
		t.Error("unknown name rendered without error")
	}
}

func TestMeta(t *testing.T) {
	all := All()
	known := map[string]bool{}
	for _, c := range meta.Categories {
		known[c.Name] = true
	}
	aliases := map[string]Name{}
	for _, n := range all {
		e, ok := meta.Icons[n]
		if !ok {
			t.Errorf("%s: no meta.json entry", n)
			continue
		}
		if !slices.Contains([]string{"lucide", "lab", "hgmx"}, e.Source) {
			t.Errorf("%s: source %q", n, e.Source)
		}
		for _, c := range e.Categories {
			if !known[c] {
				t.Errorf("%s: unknown category %q", n, c)
			}
		}
		for _, a := range e.Aliases {
			if slices.Contains(all, Name(a)) {
				t.Errorf("%s: alias %q is also an icon", n, a)
			}
			if prev, ok := aliases[a]; ok {
				t.Errorf("alias %q names both %s and %s", a, prev, n)
			}
			aliases[a] = n
		}
	}
	for n := range meta.Icons {
		if !slices.Contains(all, n) {
			t.Errorf("%s: meta.json entry without an icon", n)
		}
	}
}

func TestCategories(t *testing.T) {
	m, err := parseMeta([]byte(`{"categories":[{"name":"mail","title":"Mail"},{"name":"arrows","title":"Arrows"}],
		"icons":{"mail-open":{"source":"lucide","categories":["mail"]},
		"inbox":{"source":"lab","categories":["mail","arrows"],"aliases":["tray"]}}}`))
	if err != nil {
		t.Fatal(err)
	}
	got := m.categories([]Name{"inbox", "mail-open"})
	want := []Category{
		{Name: "mail", Title: "Mail", Icons: []Name{"inbox", "mail-open"}},
		{Name: "arrows", Title: "Arrows", Icons: []Name{"inbox"}},
	}
	if !slices.EqualFunc(got, want, func(a, b Category) bool {
		return a.Name == b.Name && a.Title == b.Title && slices.Equal(a.Icons, b.Icons)
	}) {
		t.Errorf("categories = %v, want %v", got, want)
	}
	if a := m.Icons["inbox"].Aliases; !slices.Equal(a, []string{"tray"}) {
		t.Errorf("aliases = %v", a)
	}
}
