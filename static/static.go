// Package static serves files at set-hash URLs or individual fixed URLs.
package static

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/fs"
	"maps"
	"net/http"
	"net/url"
	"path"
	"slices"
	"strings"
	"time"
)

const (
	immutable = "public, max-age=31536000, immutable"
	daily     = "public, max-age=86400"
	dev       = "dev"
)

var types = map[string]string{
	".css":         "text/css; charset=utf-8",
	".gif":         "image/gif",
	".html":        "text/html; charset=utf-8",
	".ico":         "image/x-icon",
	".jpeg":        "image/jpeg",
	".jpg":         "image/jpeg",
	".js":          "text/javascript; charset=utf-8",
	".json":        "application/json",
	".mjs":         "text/javascript; charset=utf-8",
	".png":         "image/png",
	".svg":         "image/svg+xml",
	".txt":         "text/plain; charset=utf-8",
	".webmanifest": "application/manifest+json",
	".webp":        "image/webp",
	".woff":        "font/woff",
	".woff2":       "font/woff2",
	".xml":         "application/xml",
}

// Server serves files below a root-relative URL prefix ending in a slash.
// File provides handlers for fixed URLs outside that prefix.
type Server struct {
	prefix string
	hash   string
	files  map[string]file
	live   fs.FS
}

type file struct {
	data []byte
	etag string
}

// New reads fsys once and hashes its paths and contents as one set.
func New(prefix string, fsys fs.FS) (*Server, error) {
	files := make(map[string][]byte)
	err := fs.WalkDir(fsys, ".", func(name string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		files[name], err = fs.ReadFile(fsys, name)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("static: %w", err)
	}
	return Map(prefix, files), nil
}

// Map hashes files as one set without copying their byte slices.
// Keep the slices unchanged after this call.
func Map(prefix string, files map[string][]byte) *Server {
	s := &Server{prefix: prefix, files: make(map[string]file, len(files))}
	set := sha256.New()
	for _, name := range slices.Sorted(maps.Keys(files)) {
		f := file{data: files[name], etag: etag(files[name])}
		s.files[name] = f
		fmt.Fprintf(set, "%d:%s%s", len(name), name, f.etag)
	}
	s.hash = hex.EncodeToString(set.Sum(nil)[:6])
	return s
}

// Dev reads fsys on each request and serves files with no-cache.
func Dev(prefix string, fsys fs.FS) *Server {
	return &Server{prefix: prefix, hash: dev, live: fsys}
}

// URL returns a name's escaped, hashed URL, or an error if it is missing.
// Names are unescaped paths relative to the file set's root.
func (s *Server) URL(name string) (string, error) {
	if _, ok := s.open(name); !ok {
		return "", fmt.Errorf("static: no file %q", name)
	}
	segments := strings.Split(name, "/")
	for i, seg := range segments {
		segments[i] = url.PathEscape(seg)
	}
	return s.prefix + s.hash + "/" + strings.Join(segments, "/"), nil
}

// File returns a handler for name at a fixed URL chosen by the caller.
// It panics immediately if name is missing from the file set.
// Responses cache publicly for one day, or use no-cache in dev mode.
func (s *Server) File(name string) http.Handler {
	if _, ok := s.open(name); !ok {
		panic(fmt.Sprintf("static: no file %q", name))
	}
	cache := daily
	if s.live != nil {
		cache = "no-cache"
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		s.serve(w, r, name, cache)
	})
}

// ServeHTTP expects the full URL path, without http.StripPrefix.
// A stale hash serves the current file with no-cache, not a 404.
// Current snapshot hashes get immutable caching; all others get no-cache.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	rest, ok := strings.CutPrefix(r.URL.Path, s.prefix)
	hash, name, found := strings.Cut(rest, "/")
	if !ok || !found {
		http.NotFound(w, r)
		return
	}
	cache := "no-cache"
	if hash == s.hash && s.live == nil {
		cache = immutable
	}
	s.serve(w, r, name, cache)
}

func (s *Server) serve(w http.ResponseWriter, r *http.Request, name, cache string) {
	f, ok := s.open(name)
	if !ok {
		http.NotFound(w, r)
		return
	}
	if typ, ok := types[path.Ext(name)]; ok {
		w.Header().Set("Content-Type", typ)
	}
	w.Header().Set("Cache-Control", cache)
	w.Header().Set("ETag", f.etag)
	http.ServeContent(w, r, name, time.Time{}, bytes.NewReader(f.data))
}

func (s *Server) open(name string) (file, bool) {
	if s.live == nil {
		f, ok := s.files[name]
		return f, ok
	}
	data, err := fs.ReadFile(s.live, name)
	if err != nil {
		return file{}, false
	}
	return file{data: data, etag: etag(data)}, true
}

func etag(data []byte) string {
	sum := sha256.Sum256(data)
	return `"` + hex.EncodeToString(sum[:8]) + `"`
}
