package otlptest

import (
	"bytes"
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cogikyo/hgmx/server"
	"github.com/cogikyo/hgmx/telemetry"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
)

func TestSite(t *testing.T) {
	k := &sink{t: t}
	receiver := httptest.NewServer(k)
	defer receiver.Close()
	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", receiver.URL+"/otlp?tenant=x")
	var journal bytes.Buffer
	p, err := telemetry.FromEnv(slog.NewTextHandler(&journal, nil))
	if err != nil {
		t.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("POST /items", func(w http.ResponseWriter, r *http.Request) {})
	mux.HandleFunc("GET /boom", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600")
		panic("boom")
	})
	mux.HandleFunc("GET /late", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "partial")
		panic("late")
	})
	h := server.Site{Version: "v1", InternalError: "fallback", Telemetry: p}.Handler(mux)

	serve := func(r *http.Request) *httptest.ResponseRecorder {
		w := httptest.NewRecorder()
		func() {
			defer func() { recover() }()
			h.ServeHTTP(w, r)
		}()
		return w
	}
	if w := serve(httptest.NewRequest("GET", "/healthz", nil)); w.Body.String() != "ok v1" {
		t.Errorf("healthz: %d %q", w.Code, w.Body)
	}
	serve(httptest.NewRequest("GET", "/items/7", nil))
	boom := serve(httptest.NewRequest("GET", "/boom", nil))
	if boom.Code != 500 || boom.Body.String() != "fallback" || boom.Header().Get("Content-Security-Policy") == "" ||
		boom.Header().Get("X-Content-Type-Options") != "nosniff" || boom.Header().Get("Cache-Control") != "no-store" {
		t.Errorf("panic: %d %q %v; want the fallback with security headers only", boom.Code, boom.Body, boom.Header())
	}
	serve(httptest.NewRequest("GET", "/late", nil))
	cross := httptest.NewRequest("POST", "/items", nil)
	cross.Header.Set("Sec-Fetch-Site", "cross-site")
	if w := serve(cross); w.Code != http.StatusForbidden {
		t.Errorf("cross-site POST: %d, want 403", w.Code)
	}
	if err := p.Shutdown(context.Background()); err != nil {
		t.Fatal(err)
	}

	type span struct {
		route  any
		status any
		failed bool
	}
	got := map[string]span{}
	for _, td := range k.traces {
		for _, rs := range td.ResourceSpans {
			for _, ss := range rs.ScopeSpans {
				for _, s := range ss.Spans {
					a := flat(s.Attributes)
					got[s.Name] = span{a["http.route"], a["http.response.status_code"], s.Status.GetCode() == tracepb.Status_STATUS_CODE_ERROR}
				}
			}
		}
	}
	want := map[string]span{
		"GET /items/{id}": {"/items/{id}", int64(200), false},
		"GET /boom":       {"/boom", int64(500), true},
		"GET /late":       {"/late", int64(200), true},
		"POST":            {nil, int64(403), false},
	}
	if len(got) != len(want) {
		t.Errorf("spans %v, want %v (none for /healthz)", got, want)
	}
	for name, w := range want {
		if got[name] != w {
			t.Errorf("span %q = %+v, want %+v", name, got[name], w)
		}
	}
}
