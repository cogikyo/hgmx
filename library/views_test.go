package views

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/fs"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/a-h/templ"
	"github.com/cogikyo/hgmx/scripts"
	"github.com/cogikyo/hgmx/server"
	"github.com/cogikyo/hgmx/static"
)

func TestRenderFailure(t *testing.T) {
	sub, err := fs.Sub(files, "static")
	if err != nil {
		t.Fatal(err)
	}
	st, err := static.New(assets, sub)
	if err != nil {
		t.Fatal(err)
	}
	s := &site{static: st, scripts: scripts.New(scripts.Navigation), version: "v1"}
	broken := templ.ComponentFunc(func(context.Context, io.Writer) error { return errors.New("broken") })
	page := s.page(http.StatusOK, "Broken", broken)
	var logged bytes.Buffer
	h := server.Recover(slog.New(slog.NewTextHandler(&logged, nil)), http.HandlerFunc(InternalError), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=600")
		w.Header().Set("ETag", `"v1"`)
		page.ServeHTTP(w, r)
	}))

	for _, method := range []string{"GET", "HEAD"} {
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, httptest.NewRequest(method, "/", nil))
		header := rec.Header()
		if rec.Code != http.StatusInternalServerError || header.Get("Cache-Control") != "no-store" || header.Get("ETag") != "" {
			t.Fatalf("%s: status %d, Cache-Control %q, ETag %q", method, rec.Code, header.Get("Cache-Control"), header.Get("ETag"))
		}
		if header.Get("Content-Length") != strconv.Itoa(len(internalError)) {
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
