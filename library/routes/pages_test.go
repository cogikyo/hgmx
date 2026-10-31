package routes

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/cogikyo/hgmx/server"
	"github.com/cogikyo/hgmx/telemetry"

	"github.com/cogikyo/hgmx/library/views"
)

func TestRenderFailure(t *testing.T) {
	v, err := views.New(views.Config{Version: "v1"})
	if err != nil {
		t.Fatal(err)
	}
	p := &pages{views: v}
	broken := templ.ComponentFunc(func(context.Context, io.Writer) error { return errors.New("broken") })
	page := p.page(http.StatusOK, "Broken", broken)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=600")
		w.Header().Set("ETag", `"v1"`)
		page.ServeHTTP(w, r)
	})
	var logged bytes.Buffer
	h := server.Site{
		Version:       "v1",
		InternalError: http.HandlerFunc(InternalError),
		Telemetry:     telemetry.Disabled(slog.NewTextHandler(&logged, &slog.HandlerOptions{Level: slog.LevelWarn})),
	}.Handler(mux)

	for _, method := range []string{"GET", "HEAD"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(method, "/", nil))
		header := rec.Header()
		if rec.Code != http.StatusInternalServerError || header.Get("Cache-Control") != "no-store" || header.Get("ETag") != "" {
			t.Fatalf("%s: status %d, Cache-Control %q, ETag %q", method, rec.Code, header.Get("Cache-Control"), header.Get("ETag"))
		}
		if header.Get("Content-Length") != strconv.Itoa(len(views.InternalError)) {
			t.Fatalf("%s: Content-Length %q", method, header.Get("Content-Length"))
		}
		body := rec.Body.String()
		if method == "GET" && !strings.Contains(body, `<main id="main">`) || method == "HEAD" && body != "" {
			t.Fatalf("%s: body %q", method, body)
		}
	}

	gone, cancel := context.WithCancel(context.Background())
	cancel()
	logged.Reset()
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil).WithContext(gone))
	if logged.Len() != 0 || rec.Body.Len() != 0 {
		t.Fatalf("client gone: logged %q, body %q", logged.String(), rec.Body.String())
	}
}
