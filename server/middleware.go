package server

import (
	"log/slog"
	"maps"
	"net/http"
	"runtime/debug"
	"time"
)

type response struct {
	http.ResponseWriter
	status int
	bytes  int64
}

func (w *response) WriteHeader(code int) {
	if w.status == 0 && (code >= 200 || code == http.StatusSwitchingProtocols) {
		w.status = code
	}
	w.ResponseWriter.WriteHeader(code)
}

func (w *response) Write(b []byte) (int, error) {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(b)
	w.bytes += int64(n)
	return n, err
}

func (w *response) Flush() {
	if w.status == 0 {
		w.status = http.StatusOK
	}
	http.NewResponseController(w.ResponseWriter).Flush()
}

func (w *response) Unwrap() http.ResponseWriter {
	return w.ResponseWriter
}

func Log(log *slog.Logger, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rw := &response{ResponseWriter: w}
		defer func() {
			v := recover()
			status := rw.status
			if status == 0 && v == nil {
				status = http.StatusOK
			}
			level := slog.LevelInfo
			if status >= 500 || v != nil {
				level = slog.LevelError
			}
			attrs := []slog.Attr{
				slog.String("method", r.Method),
				slog.String("pattern", r.Pattern),
				slog.String("path", r.URL.Path),
				slog.Int("status", status),
				slog.Int64("bytes", rw.bytes),
				slog.Duration("dur", time.Since(start)),
			}
			if v != nil {
				attrs = append(attrs, slog.Bool("aborted", true))
			}
			log.LogAttrs(r.Context(), level, "http", attrs...)
			if v != nil {
				panic(v)
			}
		}()
		h.ServeHTTP(rw, r)
	})
}

func Recover(log *slog.Logger, fallback, h http.Handler) http.Handler {
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
			if rw.status != 0 {
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
	"Permissions-Policy":         "accelerometer=(), browsing-topics=(), camera=(), geolocation=(), gyroscope=(), magnetometer=(), microphone=(), payment=(), usb=()",
	"Cross-Origin-Opener-Policy": "same-origin",
}

func Security(csp string, h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()
		if csp != "" {
			header.Set("Content-Security-Policy", csp)
		}
		for k, v := range headers {
			header.Set(k, v)
		}
		h.ServeHTTP(w, r)
	})
}
