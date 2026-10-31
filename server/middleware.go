package server

import (
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"net/http"
	"runtime/debug"
	"slices"
	"strings"
)

type response struct {
	http.ResponseWriter
	committed bool
}

func (w *response) WriteHeader(code int) {
	if code >= 200 || code == http.StatusSwitchingProtocols {
		w.committed = true
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *response) Write(b []byte) (int, error) {
	w.committed = true
	return w.ResponseWriter.Write(b)
}

func (w *response) FlushError() error {
	err := http.NewResponseController(w.ResponseWriter).Flush()
	if !errors.Is(err, http.ErrNotSupported) {
		w.committed = true
	}
	return err
}

func (w *response) Flush() {
	w.FlushError()
}

func (w *response) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

// recovery logs panics and serves fallback if no response has been written.
// It first restores headers to their state on entry.
// After a write or flush it aborts with http.ErrAbortHandler instead.
// An existing http.ErrAbortHandler panic passes through without logging.
func recovery(log *slog.Logger, fallback, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header().Clone()
		rw := &response{ResponseWriter: w}
		defer func() {
			v := recover()
			if v == nil {
				return
			}
			if v == http.ErrAbortHandler {
				panic(v)
			}
			log.ErrorContext(r.Context(), "panic", "method", r.Method, "path", r.URL.Path, "err", v, "stack", string(debug.Stack()))
			if rw.committed {
				panic(http.ErrAbortHandler)
			}
			clear(w.Header())
			maps.Copy(w.Header(), header)
			fallback.ServeHTTP(w, r)
		}()
		h.ServeHTTP(rw, r)
	})
}

var headers = map[string]string{
	"X-Content-Type-Options":     "nosniff",
	"Referrer-Policy":            "same-origin",
	"Cross-Origin-Opener-Policy": "same-origin",
}

var features = []string{"accelerometer", "browsing-topics", "camera", "geolocation", "gyroscope", "magnetometer", "microphone", "payment", "usb"}

// security sets default headers that h may override; the caller owns the CSP.
// An empty csp leaves Content-Security-Policy unset by this middleware.
func security(csp string, allow []string, h http.Handler) http.Handler {
	for _, name := range allow {
		if !slices.Contains(features, name) {
			panic(fmt.Sprintf("server: unknown Site.Allow feature %q", name))
		}
	}
	policy := make([]string, len(features))
	for i, name := range features {
		value := "()"
		if slices.Contains(allow, name) {
			value = "(self)"
		}
		policy[i] = name + "=" + value
	}
	permissions := strings.Join(policy, ", ")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		if csp != "" {
			header.Set("Content-Security-Policy", csp)
		}
		for k, v := range headers {
			header.Set(k, v)
		}
		header.Set("Permissions-Policy", permissions)
		h.ServeHTTP(w, r)
	})
}
