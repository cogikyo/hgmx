package telemetry

import (
	"bytes"
	"log/slog"
	"net/http"
	"net/http/httptest"
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

func TestHandlerUnwrapsLogs(t *testing.T) {
	p := offline(nil)
	var out bytes.Buffer
	h := p.Handler(p.Logs(slog.NewTextHandler(&out, nil)), http.NotFoundHandler())
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	if n := strings.Count(out.String(), "trace_id="); n != 1 {
		t.Errorf("access line has %d trace_id attrs, want 1: %s", n, out.String())
	}
	if recs := p.recs.take(); len(recs) != 0 {
		t.Errorf("access line exported as %d log records", len(recs))
	}
}
