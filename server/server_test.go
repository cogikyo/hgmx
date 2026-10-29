package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

var discard = slog.New(slog.DiscardHandler)

func freeAddr(t *testing.T) string {
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	return ln.Addr().String()
}

func serve(t *testing.T, h http.Handler) (addr string, cancel context.CancelFunc, done <-chan error) {
	addr = freeAddr(t)
	ctx, cancel := context.WithCancel(context.Background())
	errc := make(chan error, 1)
	t.Cleanup(cancel)
	go func() { errc <- Run(ctx, New(addr, h, discard)) }()
	for range 100 {
		select {
		case err := <-errc:
			t.Fatalf("Run returned before serving: %v", err)
		default:
		}
		if c, err := net.Dial("tcp", addr); err == nil {
			c.Close()
			return addr, cancel, errc
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("server did not start")
	return
}

func TestRun(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	blocked := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		close(entered)
		<-release
		io.WriteString(w, "done")
	})

	t.Run("drain", func(t *testing.T) {
		addr, cancel, done := serve(t, blocked)
		got := make(chan string, 1)
		go func() {
			res, err := http.Get("http://" + addr)
			if err != nil {
				got <- err.Error()
				return
			}
			b, _ := io.ReadAll(res.Body)
			res.Body.Close()
			got <- string(b)
		}()
		<-entered
		cancel()
		select {
		case err := <-done:
			t.Fatalf("Run returned %v before the handler finished", err)
		case <-time.After(100 * time.Millisecond):
		}
		close(release)
		if err := <-done; err != nil {
			t.Fatalf("Run = %v, want nil", err)
		}
		if body := <-got; body != "done" {
			t.Fatalf("drained response = %q, want done", body)
		}
	})

	t.Run("deadline", func(t *testing.T) {
		defer func(d time.Duration) { drainTimeout = d }(drainTimeout)
		drainTimeout = 50 * time.Millisecond
		entered, stuck := make(chan struct{}), make(chan struct{})
		defer close(stuck)
		addr, cancel, done := serve(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			close(entered)
			<-stuck
		}))
		client := make(chan error, 1)
		go func() {
			res, err := http.Get("http://" + addr)
			if err == nil {
				res.Body.Close()
			}
			client <- err
		}()
		<-entered
		cancel()
		select {
		case err := <-done:
			if !errors.Is(err, context.DeadlineExceeded) {
				t.Fatalf("Run = %v, want deadline exceeded", err)
			}
		case <-time.After(5 * time.Second):
			t.Fatal("Run did not return after the drain deadline")
		}
		select {
		case err := <-client:
			if err == nil {
				t.Fatal("request succeeded, want its connection closed at the deadline")
			}
		case <-time.After(5 * time.Second):
			t.Fatal("request still open after the drain deadline")
		}
	})

	t.Run("bind", func(t *testing.T) {
		ln, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			t.Fatal(err)
		}
		defer ln.Close()
		if err := Run(context.Background(), New(ln.Addr().String(), blocked, discard)); err == nil {
			t.Fatal("Run on a bound address returned nil")
		}
	})
}

type entry struct {
	Pattern string
	Status  int
	Bytes   int64
	Aborted bool
}

func observe(t *testing.T, h http.Handler, req *http.Request) (*httptest.ResponseRecorder, entry) {
	t.Helper()
	var buf bytes.Buffer
	log := slog.New(slog.NewJSONHandler(&buf, nil))
	rec := httptest.NewRecorder()
	func() {
		defer func() {
			if v := recover(); v != nil && v != http.ErrAbortHandler {
				t.Fatalf("panic %v escaped", v)
			}
		}()
		Log(log, h).ServeHTTP(rec, req)
	}()
	var e entry
	if err := json.Unmarshal(buf.Bytes(), &e); err != nil {
		t.Fatalf("observer log %q: %v", buf.String(), err)
	}
	return rec, e
}

func TestLog(t *testing.T) {
	for _, tc := range []struct {
		name string
		h    http.HandlerFunc
		want entry
	}{
		{"implicit", func(w http.ResponseWriter, r *http.Request) {}, entry{Status: 200}},
		{"write", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "abc") }, entry{Status: 200, Bytes: 3}},
		{"explicit", func(w http.ResponseWriter, r *http.Request) { w.WriteHeader(404) }, entry{Status: 404}},
		{"first final", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusEarlyHints)
			w.WriteHeader(201)
			w.WriteHeader(500)
		}, entry{Status: 201}},
		{"flush", func(w http.ResponseWriter, r *http.Request) {
			if err := http.NewResponseController(w).Flush(); err != nil {
				panic(err)
			}
		}, entry{Status: 200}},
		{"aborted", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "ab")
			panic(http.ErrAbortHandler)
		}, entry{Status: 200, Bytes: 2, Aborted: true}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, got := observe(t, tc.h, httptest.NewRequest("GET", "/", nil))
			if got != tc.want {
				t.Fatalf("observed %+v, want %+v", got, tc.want)
			}
		})
	}

	t.Run("flush error", func(t *testing.T) {
		wrappers := map[string]func(http.Handler) http.Handler{
			"Log":     func(h http.Handler) http.Handler { return Log(discard, h) },
			"Recover": func(h http.Handler) http.Handler { return Recover(discard, http.NotFoundHandler(), h) },
		}
		for name, wrap := range wrappers {
			for _, tc := range []struct {
				w    http.ResponseWriter
				want error
			}{
				{failing{httptest.NewRecorder()}, errGone},
				{struct{ http.ResponseWriter }{httptest.NewRecorder()}, http.ErrNotSupported},
			} {
				var got error
				wrap(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					got = http.NewResponseController(w).Flush()
				})).ServeHTTP(tc.w, httptest.NewRequest("GET", "/", nil))
				if !errors.Is(got, tc.want) {
					t.Errorf("%s: Flush = %v, want %v", name, got, tc.want)
				}
			}
		}
	})
}

var errGone = errors.New("client gone")

type failing struct{ *httptest.ResponseRecorder }

func (failing) FlushError() error { return errGone }

func TestChain(t *testing.T) {
	const csp = "default-src 'self'"
	site := http.NewServeMux()
	site.HandleFunc("GET /items/{id}", func(w http.ResponseWriter, r *http.Request) {})
	site.HandleFunc("POST /items", func(w http.ResponseWriter, r *http.Request) {})
	site.HandleFunc("GET /boom", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Header().Set("ETag", `"x"`)
		panic("boom")
	})
	site.HandleFunc("GET /late", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "partial")
		panic("late")
	})
	fallback := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "<!doctype html>error")
	})
	chain := Security(csp, Recover(discard, fallback, http.NewCrossOriginProtection().Handler(site)))

	cross := httptest.NewRequest("POST", "/items", nil)
	cross.Header.Set("Sec-Fetch-Site", "cross-site")

	for _, tc := range []struct {
		name   string
		req    *http.Request
		want   entry
		body   string
		absent []string
	}{
		{"pattern", httptest.NewRequest("GET", "/items/7", nil), entry{Pattern: "GET /items/{id}", Status: 200}, "", nil},
		{"recovered", httptest.NewRequest("GET", "/boom", nil), entry{Pattern: "GET /boom", Status: 500, Bytes: 20}, "<!doctype html>error", []string{"Cache-Control", "ETag"}},
		{"committed", httptest.NewRequest("GET", "/late", nil), entry{Pattern: "GET /late", Status: 200, Bytes: 7, Aborted: true}, "partial", nil},
		{"cross-origin", cross, entry{Status: 403, Bytes: 57}, "", nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			rec, got := observe(t, chain, tc.req)
			if got != tc.want {
				t.Fatalf("observed %+v, want %+v", got, tc.want)
			}
			if tc.body != "" && rec.Body.String() != tc.body {
				t.Fatalf("body %q, want %q", rec.Body.String(), tc.body)
			}
			if rec.Header().Get("Content-Security-Policy") != csp || rec.Header().Get("X-Content-Type-Options") != "nosniff" {
				t.Fatalf("security headers missing: %v", rec.Header())
			}
			for _, k := range tc.absent {
				if v := rec.Header().Get(k); v != "" {
					t.Fatalf("%s = %q leaked from the failed handler", k, v)
				}
			}
		})
	}
}

func TestPageVersion(t *testing.T) {
	h := PageVersion("v2", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	}))
	for _, tc := range []struct {
		method, header string
		redirect       bool
	}{
		{"GET", "v1", true},
		{"HEAD", "v1", true},
		{"GET", "v2", false},
		{"GET", "", false},
		{"POST", "v1", false},
	} {
		req := httptest.NewRequest(tc.method, "/docs?q=a", nil)
		if tc.header != "" {
			req.Header.Set(versionHeader, tc.header)
		}
		rec := httptest.NewRecorder()
		h.ServeHTTP(rec, req)
		redirected := rec.Header().Get("HX-Redirect") == "/docs?q=a" && rec.Code == http.StatusNoContent
		if redirected != tc.redirect || !redirected && rec.Code != http.StatusTeapot {
			t.Errorf("%s with version %q: status %d, HX-Redirect %q", tc.method, tc.header, rec.Code, rec.Header().Get("HX-Redirect"))
		}
	}
}
