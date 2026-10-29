package telemetry

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"testing"
	"testing/slogtest"
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
	h := p.Logs(slog.DiscardHandler)
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

func TestLogsTraceIDsAtTopLevel(t *testing.T) {
	p := offline(nil)
	var out bytes.Buffer
	l := slog.New(p.Logs(slog.NewTextHandler(&out, nil))).With("a", 1).WithGroup("g").With("b", 2)
	ctx, s := p.Start(context.Background(), "work")
	l.InfoContext(ctx, "hi", "c", 3)
	line := out.String()
	want := "a=1 trace_id=" + s.TraceID() + " span_id=" + s.spanHex() + " g.b=2 g.c=3"
	if !strings.Contains(line, want) {
		t.Errorf("local line %q lacks %q", line, want)
	}
}
