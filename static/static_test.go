package static

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func get(t *testing.T, h http.Handler, url string, header ...string) *httptest.ResponseRecorder {
	t.Helper()
	r := httptest.NewRequest(http.MethodGet, url, nil)
	for i := 0; i+1 < len(header); i += 2 {
		r.Header.Set(header[i], header[i+1])
	}
	w := httptest.NewRecorder()
	h.ServeHTTP(w, r)
	return w
}

func TestCaching(t *testing.T) {
	s := Map("/static/", map[string][]byte{"fonts/Sans Serif.woff2": []byte("wOF2")})
	url, err := s.URL("fonts/Sans Serif.woff2")
	if err != nil {
		t.Fatal(err)
	}

	current := get(t, s, url)
	if current.Code != http.StatusOK || !strings.Contains(current.Header().Get("Cache-Control"), "immutable") {
		t.Fatalf("current hash: %d %q, want 200 immutable", current.Code, current.Header().Get("Cache-Control"))
	}
	stale := get(t, s, "/static/000000000000/fonts/Sans%20Serif.woff2")
	if stale.Code != http.StatusOK || stale.Header().Get("Cache-Control") != "no-cache" {
		t.Fatalf("stale hash: %d %q, want 200 no-cache", stale.Code, stale.Header().Get("Cache-Control"))
	}
	if again := get(t, s, url, "If-None-Match", current.Header().Get("ETag")); again.Code != http.StatusNotModified {
		t.Fatalf("If-None-Match: %d, want 304", again.Code)
	}
}

func TestEditMovesEveryURL(t *testing.T) {
	before := Map("/static/", map[string][]byte{"css/main.css": []byte("@import 'a.css';"), "css/a.css": []byte("a{}")})
	after := Map("/static/", map[string][]byte{"css/main.css": []byte("@import 'a.css';"), "css/a.css": []byte("a{color:red}")})
	old, _ := before.URL("css/main.css")
	now, _ := after.URL("css/main.css")
	if old == now {
		t.Fatalf("editing css/a.css left css/main.css at %s, so browsers keep the cached copy that imports the old a.css", now)
	}
}

func TestDevServesEdits(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "app.js")
	if err := os.WriteFile(file, []byte("one"), 0o644); err != nil {
		t.Fatal(err)
	}
	s := Dev("/static/", os.DirFS(dir))
	url, err := s.URL("app.js")
	if err != nil {
		t.Fatal(err)
	}
	first := get(t, s, url)
	if err := os.WriteFile(file, []byte("two"), 0o644); err != nil {
		t.Fatal(err)
	}
	second := get(t, s, url, "If-None-Match", first.Header().Get("ETag"))
	if second.Code != http.StatusOK || second.Body.String() != "two" || second.Header().Get("Cache-Control") != "no-cache" {
		t.Fatalf("after edit: %d %q %q, want 200 \"two\" no-cache", second.Code, second.Body.String(), second.Header().Get("Cache-Control"))
	}
}
