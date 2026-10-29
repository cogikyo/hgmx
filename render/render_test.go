package render

import (
	"cmp"
	"context"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/a-h/templ"
)

const page = "<!doctype html><title>t</title>"

var component = templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
	_, err := io.WriteString(w, page)
	return err
})

func TestHTML(t *testing.T) {
	first := httptest.NewRecorder()
	if err := HTML(first, httptest.NewRequest("GET", "/", nil), http.StatusOK, component); err != nil {
		t.Fatal(err)
	}
	etag := first.Header().Get("ETag")
	if etag == "" || first.Body.String() != page {
		t.Fatalf("first response: ETag %q, body %q", etag, first.Body.String())
	}

	for _, tc := range []struct {
		name, method, inm, cache string
		status, code             int
		body                     string
	}{
		{"get", "GET", "", "", 200, 200, page},
		{"head", "HEAD", "", "", 200, 200, ""},
		{"strong", "GET", etag, "", 200, 304, ""},
		{"weak", "GET", "W/" + etag, "", 200, 304, ""},
		{"list", "GET", `"other", ` + etag, "", 200, 304, ""},
		{"miss", "GET", `"other"`, "", 200, 200, page},
		{"wildcard", "HEAD", "*", "", 200, 304, ""},
		{"unprocessable", "GET", "*", "", 422, 422, page},
		{"not found head", "HEAD", "*", "", 404, 404, ""},
		{"mutation", "POST", "*", "", 200, 200, page},
		{"no-store", "GET", "", "no-store", 200, 200, page},
	} {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/", nil)
			if tc.inm != "" {
				req.Header.Set("If-None-Match", tc.inm)
			}
			rec := httptest.NewRecorder()
			if tc.cache != "" {
				rec.Header().Set("Cache-Control", tc.cache)
			}
			if err := HTML(rec, req, tc.status, component); err != nil {
				t.Fatal(err)
			}
			if rec.Code != tc.code || rec.Body.String() != tc.body {
				t.Fatalf("got %d %q, want %d %q", rec.Code, rec.Body.String(), tc.code, tc.body)
			}
			if want := cmp.Or(tc.cache, "no-cache"); rec.Header().Get("Cache-Control") != want {
				t.Fatalf("Cache-Control %q, want %q", rec.Header().Get("Cache-Control"), want)
			}
			if tc.method == "HEAD" && tc.code == 200 && rec.Header().Get("Content-Length") != "31" {
				t.Fatalf("HEAD Content-Length %q, want 31", rec.Header().Get("Content-Length"))
			}
		})
	}
}

func TestHTMLFailure(t *testing.T) {
	broken := templ.ComponentFunc(func(ctx context.Context, w io.Writer) error {
		io.WriteString(w, "<!doctype html><p>half")
		return errors.New("broken")
	})
	rec := httptest.NewRecorder()
	if err := HTML(rec, httptest.NewRequest("GET", "/", nil), http.StatusOK, broken); err == nil {
		t.Fatal("HTML returned nil for a failing component")
	}
	if len(rec.Header()) != 0 || rec.Body.Len() != 0 || rec.Flushed {
		t.Fatalf("failed render leaked: headers %v, body %q", rec.Header(), rec.Body.String())
	}
}
