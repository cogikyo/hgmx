package routes

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
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
	mux := http.NewServeMux()
	mux.Handle("/", p.page(http.StatusOK, "Broken", broken))
	var logged bytes.Buffer
	h := server.Site{
		Version:   "v1",
		Telemetry: telemetry.Disabled(slog.NewTextHandler(&logged, &slog.HandlerOptions{Level: slog.LevelWarn})),
	}.Handler(mux)

	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil))
	if rec.Code != http.StatusInternalServerError || !strings.Contains(logged.String(), "broken") {
		t.Fatalf("status %d, logged %q", rec.Code, logged.String())
	}

	gone, cancel := context.WithCancel(context.Background())
	cancel()
	logged.Reset()
	rec = httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest("GET", "/", nil).WithContext(gone))
	if logged.Len() != 0 || rec.Body.Len() != 0 {
		t.Fatalf("client gone: logged %q, body %q", logged.String(), rec.Body.String())
	}
}

func TestNew(t *testing.T) {
	v, err := views.New(views.Config{Version: "v1"})
	if err != nil {
		t.Fatal(err)
	}
	New(nil, v)
}
