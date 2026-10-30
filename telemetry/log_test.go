package telemetry

import (
	"log/slog"
	"strings"
	"testing"
	"testing/slogtest"
	"unicode/utf8"
)

func tree(attrs []slog.Attr) map[string]any {
	m := map[string]any{}
	for _, a := range attrs {
		if a.Value.Kind() == slog.KindGroup {
			m[a.Key] = tree(a.Value.Group())
		} else {
			m[a.Key] = a.Value.Any()
		}
	}
	return m
}

func TestLogsSlogtest(t *testing.T) {
	p := offline(nil)
	h := p.Logger().Handler()
	slogtest.Run(t, func(*testing.T) slog.Handler { return h }, func(t *testing.T) map[string]any {
		recs := p.recs.take()
		if len(recs) != 1 {
			t.Fatalf("%d exported records", len(recs))
		}
		r := recs[0]
		m := tree(r.attrs)
		m[slog.LevelKey] = r.level
		m[slog.MessageKey] = r.msg
		if r.time != 0 {
			m[slog.TimeKey] = r.time
		}
		return m
	})
}

func TestClipRepairsUTF8(t *testing.T) {
	got := clip(strings.Repeat("\xffa", 1000))
	if len(got) > maxString || !utf8.ValidString(got) {
		t.Fatalf("clip len=%d valid=%v", len(got), utf8.ValidString(got))
	}
	if got := clip(strings.Repeat("\x80", 2000)); got != "\uFFFD" {
		t.Fatalf("clip continuation bytes = %q", got)
	}
}
