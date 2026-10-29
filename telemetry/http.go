package telemetry

import (
	"cmp"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type recorder struct {
	http.ResponseWriter
	status int
	size   int64
}

func (r *recorder) WriteHeader(status int) {
	if r.status == 0 && status >= 200 {
		r.status = status
	}
	r.ResponseWriter.WriteHeader(status)
}

func (r *recorder) Write(b []byte) (int, error) {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	n, err := r.ResponseWriter.Write(b)
	r.size += int64(n)
	return n, err
}

func (r *recorder) Flush() {
	if r.status == 0 {
		r.status = http.StatusOK
	}
	http.NewResponseController(r.ResponseWriter).Flush()
}

func (r *recorder) Unwrap() http.ResponseWriter { return r.ResponseWriter }

// Handler records requests and writes access logs to the raw local handler.
// A handler from [Pipeline.Logs] is unwrapped to avoid exporting access logs.
// The first local handler given to Logs or Handler gets export-failure warnings.
//
// Wrap the mux directly, or only middleware passing the same *http.Request.
// An intervening r.WithContext loses http.route.
// Keep /healthz and [Beacon] outside this handler.
func (p *Pipeline) Handler(local slog.Handler, next http.Handler) http.Handler {
	local = p.local(local)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &recorder{ResponseWriter: w}
		s := p.root("", server)
		if s.recording() {
			if l, ok := parseTraceparent(r.Header.Get("traceparent")); ok {
				s.links = []link{l}
			}
			w.Header().Add("Server-Timing", `traceparent;desc="`+s.traceparent()+`"`)
		}
		if s != nil {
			r = r.WithContext(context.WithValue(r.Context(), spanKey{}, s))
			p.meters.active.add(1)
		}
		served := false
		defer func() {
			status := rec.status
			switch {
			case !served:
				status = http.StatusInternalServerError
			case status == 0:
				status = http.StatusOK
			}
			d := time.Since(start)
			if s != nil {
				p.meters.active.add(-1)
				p.observe(s, r, status, rec.size, d)
			}
			access(local, r, s, status, d)
		}()
		next.ServeHTTP(rec, r)
		served = true
	})
}

func (p *Pipeline) observe(s *Span, r *http.Request, status int, size int64, d time.Duration) {
	method := normalize(r.Method)
	route := route(r.Pattern)
	attrs := make([]slog.Attr, 0, 6)
	attrs = append(attrs,
		slog.String("http.request.method", method),
		slog.Int("http.response.status_code", status),
	)
	if route != "" {
		attrs = append(attrs, slog.String("http.route", route))
	}
	p.meters.duration.record(d.Seconds(), attrs...)
	if !s.recording() {
		return
	}
	s.name = method
	if route != "" {
		s.name = method + " " + route
	}
	attrs = append(attrs, slog.String("url.path", r.URL.Path), slog.Int64("http.response.body.size", size))
	if method != r.Method {
		attrs = append(attrs, slog.String("http.request.method_original", r.Method))
	}
	s.Set(attrs...)
	if status >= 500 {
		s.code = 2
	}
	s.End()
}

func access(local slog.Handler, r *http.Request, s *Span, status int, d time.Duration) {
	level := slog.LevelInfo
	if status >= 500 {
		level = slog.LevelError
	}
	if local == nil || !local.Enabled(r.Context(), level) {
		return
	}
	rec := slog.NewRecord(time.Now(), level, "request", 0)
	rec.AddAttrs(
		slog.String("method", r.Method),
		slog.Int("status", status),
		slog.String("path", r.URL.Path),
		slog.Duration("dur", d.Truncate(time.Millisecond)),
		slog.String("host", r.Host),
	)
	if s.recording() {
		rec.AddAttrs(slog.String("trace_id", s.TraceID()), slog.String("span_id", s.spanHex()))
	}
	local.Handle(r.Context(), rec)
}

func route(pattern string) string {
	if i := strings.IndexByte(pattern, ' '); i >= 0 {
		pattern = pattern[i+1:]
	}
	if i := strings.IndexByte(pattern, '/'); i > 0 {
		pattern = pattern[i:]
	}
	return strings.TrimSuffix(pattern, "{$}")
}

var methods = map[string]string{
	http.MethodGet: http.MethodGet, http.MethodHead: http.MethodHead, http.MethodPost: http.MethodPost,
	http.MethodPut: http.MethodPut, http.MethodPatch: http.MethodPatch, http.MethodDelete: http.MethodDelete,
	http.MethodConnect: http.MethodConnect, http.MethodOptions: http.MethodOptions, http.MethodTrace: http.MethodTrace,
}

func normalize(method string) string {
	if m, ok := methods[method]; ok {
		return m
	}
	return "_OTHER"
}

type transport struct {
	p    *Pipeline
	base http.RoundTripper
}

// Transport wraps base with client spans; nil base means http.DefaultTransport.
// Build one http.Client with it and pass the parent span's request context.
// Without a span in the context it starts a rate-limited root client span.
// Spans end on response headers, before the body is read.
// The wrapper removes traceparent and tracestate before sending upstream.
// A nil pipeline returns base unchanged after applying the default.
func (p *Pipeline) Transport(base http.RoundTripper) http.RoundTripper {
	if base == nil {
		base = http.DefaultTransport
	}
	if p == nil {
		return base
	}
	return transport{p, base}
}

func (t transport) RoundTrip(req *http.Request) (*http.Response, error) {
	req = req.Clone(req.Context())
	req.Header.Del("traceparent")
	req.Header.Del("tracestate")
	method := normalize(cmp.Or(req.Method, http.MethodGet))
	s := t.p.begin(req.Context(), method, client)
	if !s.recording() {
		return t.base.RoundTrip(req)
	}
	port := req.URL.Port()
	if port == "" {
		port = "80"
		if req.URL.Scheme == "https" {
			port = "443"
		}
	}
	n, _ := strconv.Atoi(port)
	full := url.URL{Scheme: req.URL.Scheme, Host: req.URL.Host, Path: req.URL.Path}
	s.Set(
		slog.String("http.request.method", method),
		slog.String("server.address", req.URL.Hostname()),
		slog.Int("server.port", n),
		slog.String("url.full", full.String()),
	)
	resp, err := t.base.RoundTrip(req)
	switch {
	case err != nil:
		s.Set(slog.String("error.type", fmt.Sprintf("%T", err)))
		s.Fail(err)
	case resp.StatusCode >= 400:
		s.Set(slog.Int("http.response.status_code", resp.StatusCode), slog.String("error.type", strconv.Itoa(resp.StatusCode)))
		s.code = 2
	default:
		s.Set(slog.Int("http.response.status_code", resp.StatusCode))
	}
	s.End()
	return resp, err
}
