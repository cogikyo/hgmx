package server

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cogikyo/hgmx/telemetry"
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

func start(t *testing.T, h http.Handler, after ...func(context.Context) error) (srv *http.Server, cancel context.CancelFunc, done <-chan error) {
	srv = New(freeAddr(t), h, discard)
	ctx, cancel := context.WithCancel(context.Background())
	errc := make(chan error, 1)
	t.Cleanup(cancel)
	go func() { errc <- Run(ctx, srv, after...) }()
	for range 100 {
		select {
		case err := <-errc:
			t.Fatalf("Run returned before serving: %v", err)
		default:
		}
		if c, err := net.Dial("tcp", srv.Addr); err == nil {
			c.Close()
			return srv, cancel, errc
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatal("server did not start")
	return
}

func get(addr string) <-chan string {
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
	return got
}

func TestRun(t *testing.T) {
	t.Run("drain", func(t *testing.T) {
		entered, release := make(chan struct{}), make(chan struct{})
		srv, cancel, done := start(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			close(entered)
			<-release
			io.WriteString(w, "done")
		}))
		got := get(srv.Addr)
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

	t.Run("external shutdown", func(t *testing.T) {
		entered, release := make(chan struct{}), make(chan struct{})
		var finished atomic.Bool
		srv, _, done := start(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			close(entered)
			<-release
			io.WriteString(w, "done")
			finished.Store(true)
		}), func(context.Context) error {
			if !finished.Load() {
				t.Error("hook ran before the in-flight handler finished")
			}
			return nil
		})
		got := get(srv.Addr)
		<-entered
		go srv.Shutdown(context.Background())
		time.Sleep(100 * time.Millisecond)
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
		hook := make(chan error, 1)
		srv, cancel, done := start(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			close(entered)
			<-stuck
		}), func(ctx context.Context) error {
			hook <- ctx.Err()
			return nil
		})
		got := get(srv.Addr)
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
		if err := <-hook; err != nil {
			t.Fatalf("hook context after a full drain: %v, want live", err)
		}
		select {
		case body := <-got:
			if body == "" {
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
		errHook := errors.New("hook")
		calls := 0
		err = Run(context.Background(), New(ln.Addr().String(), http.NotFoundHandler(), discard), func(ctx context.Context) error {
			calls++
			if ctx.Err() != nil {
				t.Errorf("hook context: %v, want live", ctx.Err())
			}
			return errHook
		})
		var bind *net.OpError
		if calls != 1 || !errors.As(err, &bind) || !errors.Is(err, errHook) {
			t.Fatalf("Run on a bound address = %v after %d hook calls; want the bind and hook errors after one call", err, calls)
		}
	})
}

func TestHandlerAccess(t *testing.T) {
	var log bytes.Buffer
	mux := http.NewServeMux()
	mux.HandleFunc("GET /a/{id}", func(w http.ResponseWriter, r *http.Request) {})
	h := Site{InternalError: http.NotFoundHandler(), Telemetry: telemetry.Disabled(slog.NewTextHandler(&log, nil))}.Handler(mux)
	for _, tc := range []struct{ target, want string }{
		{"//a/7", `pattern="GET /a/{id}" path=//a/7 status=307`},
		{"/a/7", `pattern="GET /a/{id}" path=/a/7 status=200`},
	} {
		log.Reset()
		h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", tc.target, nil))
		if !strings.Contains(log.String(), tc.want) {
			t.Errorf("GET %s access line %q, want %s", tc.target, log.String(), tc.want)
		}
	}
}

var errGone = errors.New("client gone")

type failing struct{ *httptest.ResponseRecorder }

func (failing) FlushError() error { return errGone }

func TestRecoveryFlush(t *testing.T) {
	for _, tc := range []struct {
		w    http.ResponseWriter
		want error
	}{
		{failing{httptest.NewRecorder()}, errGone},
		{struct{ http.ResponseWriter }{httptest.NewRecorder()}, http.ErrNotSupported},
	} {
		var got error
		recovery(discard, http.NotFoundHandler(), http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			got = http.NewResponseController(w).Flush()
		})).ServeHTTP(tc.w, httptest.NewRequest("GET", "/", nil))
		if !errors.Is(got, tc.want) {
			t.Errorf("Flush = %v, want %v", got, tc.want)
		}
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
