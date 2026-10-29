package telemetry

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func recoverer(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if v := recover(); v != nil {
				http.Error(w, "internal", http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func headers(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Content-Type-Options", "nosniff")
		next.ServeHTTP(w, r)
	})
}

type ctxKey struct{}

func rewrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKey{}, 1)))
	})
}

func counts(p *Pipeline) map[*series]uint64 {
	d := p.meters.duration
	d.mu.Lock()
	defer d.mu.Unlock()
	out := map[*series]uint64{}
	for _, s := range d.series {
		out[s] = s.count
	}
	return out
}

func attr(attrs []slog.Attr, key string) string {
	for _, a := range attrs {
		if a.Key == key {
			return a.Value.String()
		}
	}
	return ""
}

func routes(t *testing.T, chain func(http.Handler) http.Handler, path string) (span, point string) {
	t.Helper()
	p := offline(nil)
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+path+"/{id}", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "ok") })
	h := p.Handler(chain(mux))

	before := counts(p)
	w := httptest.NewRecorder()
	h.ServeHTTP(w, httptest.NewRequest("GET", path+"/7", nil))
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}

	spans := p.spans.take()
	if len(spans) != 1 {
		t.Fatalf("%d spans", len(spans))
	}
	var points []*series
	for s, n := range counts(p) {
		if n > before[s] {
			points = append(points, s)
		}
	}
	if len(points) != 1 {
		t.Fatalf("%d duration points changed", len(points))
	}
	return attr(spans[0].attrs, "http.route"), attr(points[0].attrs, "http.route")
}

func TestRouteThroughChain(t *testing.T) {
	cop := http.NewCrossOriginProtection()
	chain := func(mux http.Handler) http.Handler { return recoverer(headers(cop.Handler(mux))) }
	span, point := routes(t, chain, "/chain")
	if span != "/chain/{id}" || point != "/chain/{id}" {
		t.Errorf("span route %q, point route %q; want /chain/{id} on both", span, point)
	}
}

func TestRouteLostByWithContext(t *testing.T) {
	span, point := routes(t, rewrap, "/rewrap")
	if span != "" || point != "" {
		t.Errorf("span route %q, point route %q; want the route hidden by r.WithContext", span, point)
	}
}
