package telemetry

import (
	"bytes"
	"cmp"
	"compress/gzip"
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"unicode/utf8"
)

const (
	scope = "github.com/cogikyo/hgmx/telemetry"

	maxItems     = 2048
	batch        = 512
	flushEvery   = 5 * time.Second
	collectEvery = 60 * time.Second
	timeout      = 5 * time.Second
	minRetry     = 5 * time.Second
	maxRetry     = 30 * time.Second
	rootsPerSec  = 200
)

var version = func() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	if info.Main.Path == "github.com/cogikyo/hgmx" {
		return info.Main.Version
	}
	for _, d := range info.Deps {
		if d.Path == "github.com/cogikyo/hgmx" {
			return d.Version
		}
	}
	return ""
}()

type signal int

const (
	sigTraces signal = iota
	sigLogs
	sigMetrics
)

var signals = [...]string{"spans", "logs", "metrics"}

type pending struct {
	body []byte
	n    int
	due  time.Time
}

type Pipeline struct {
	urls   [3]string
	res    []byte
	client *http.Client
	start  int64
	meters meters

	spans   *buffer[*Span]
	recs    *buffer[record]
	roots   *bucket
	retry   [3]*pending
	failing [3]bool
	paths   [3]string
	warn    atomic.Pointer[slog.Handler]

	closed atomic.Bool
	kick   chan struct{}
	stop   chan struct{}
	loop   chan struct{}
	done   chan struct{}
	err    error

	ctx    context.Context
	cancel context.CancelFunc
}

func (p *Pipeline) off() bool { return p == nil || p.closed.Load() }

type buffer[T any] struct {
	mu     sync.Mutex
	items  []T
	closed *atomic.Bool
	kick   chan struct{}
	drop   func(int)
}

func (b *buffer[T]) add(v T) {
	b.mu.Lock()
	if b.closed.Load() {
		b.mu.Unlock()
		return
	}
	if len(b.items) >= maxItems {
		b.mu.Unlock()
		b.drop(1)
		return
	}
	b.items = append(b.items, v)
	full := len(b.items) == batch
	b.mu.Unlock()
	if full {
		select {
		case b.kick <- struct{}{}:
		default:
		}
	}
}

func (b *buffer[T]) take() []T {
	b.mu.Lock()
	defer b.mu.Unlock()
	items := b.items
	b.items = nil
	return items
}

func FromEnv() (*Pipeline, error) {
	if strings.EqualFold(strings.TrimSpace(os.Getenv("OTEL_SDK_DISABLED")), "true") {
		return nil, nil
	}
	base := strings.TrimSpace(os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT"))
	if base == "" {
		return nil, nil
	}
	u, err := url.Parse(base)
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" {
		return nil, fmt.Errorf("telemetry: OTEL_EXPORTER_OTLP_ENDPOINT %q is not an http(s) URL", base)
	}
	if v := os.Getenv("OTEL_EXPORTER_OTLP_PROTOCOL"); v != "" && v != "http/protobuf" {
		return nil, fmt.Errorf("telemetry: OTEL_EXPORTER_OTLP_PROTOCOL %q is unsupported; use http/protobuf", v)
	}
	if v := strings.ToLower(os.Getenv("OTEL_EXPORTER_OTLP_METRICS_TEMPORALITY_PREFERENCE")); v != "" && v != "cumulative" {
		return nil, fmt.Errorf("telemetry: OTEL_EXPORTER_OTLP_METRICS_TEMPORALITY_PREFERENCE %q is unsupported; use cumulative", v)
	}
	res, err := resource(os.Getenv("OTEL_SERVICE_NAME"), os.Getenv("OTEL_RESOURCE_ATTRIBUTES"))
	if err != nil {
		return nil, err
	}
	p := newPipeline(u, res, &http.Transport{MaxIdleConns: 2, IdleConnTimeout: time.Minute})
	go p.run()
	return p, nil
}

func newPipeline(u *url.URL, res []slog.Attr, rt http.RoundTripper) *Pipeline {
	p := &Pipeline{
		res: encodeResource(res),
		client: &http.Client{
			Timeout:       timeout,
			Transport:     rt,
			CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
		},
		start:  time.Now().UnixNano(),
		meters: newMeters(),
		roots:  newBucket(rootsPerSec),
		kick:   make(chan struct{}, 1),
		stop:   make(chan struct{}),
		loop:   make(chan struct{}),
		done:   make(chan struct{}),
	}
	for i, path := range [3]string{"v1/traces", "v1/logs", "v1/metrics"} {
		p.urls[i] = u.JoinPath(path).String()
		p.paths[i] = u.JoinPath(path).Path
	}
	p.spans = &buffer[*Span]{closed: &p.closed, kick: p.kick, drop: func(n int) { p.drop(sigTraces, n) }}
	p.recs = &buffer[record]{closed: &p.closed, kick: p.kick, drop: func(n int) { p.drop(sigLogs, n) }}
	p.ctx, p.cancel = context.WithCancel(context.Background())
	return p
}

func resource(name, env string) ([]slog.Attr, error) {
	if !utf8.ValidString(name) {
		return nil, fmt.Errorf("telemetry: OTEL_SERVICE_NAME is not valid UTF-8")
	}
	var attrs []slog.Attr
	for kv := range strings.SplitSeq(env, ",") {
		if strings.TrimSpace(kv) == "" {
			continue
		}
		k, v, ok := strings.Cut(kv, "=")
		k = strings.TrimSpace(k)
		if !ok || k == "" {
			return nil, fmt.Errorf("telemetry: OTEL_RESOURCE_ATTRIBUTES entry %q is not key=value", kv)
		}
		v, err := url.PathUnescape(strings.TrimSpace(v))
		if err != nil {
			return nil, fmt.Errorf("telemetry: OTEL_RESOURCE_ATTRIBUTES entry %q: %w", kv, err)
		}
		if !utf8.ValidString(k) || !utf8.ValidString(v) {
			return nil, fmt.Errorf("telemetry: OTEL_RESOURCE_ATTRIBUTES entry %q is not valid UTF-8", kv)
		}
		if k == "service.name" {
			name = cmp.Or(strings.TrimSpace(name), v)
			continue
		}
		attrs = append(attrs, slog.String(k, v))
	}
	name = cmp.Or(strings.TrimSpace(name), "unknown_service:"+filepath.Base(os.Args[0]))
	return append(attrs,
		slog.String("service.name", name),
		slog.String("telemetry.sdk.name", "hgmx"),
		slog.String("telemetry.sdk.language", "go"),
		slog.String("telemetry.sdk.version", version),
	), nil
}

func (p *Pipeline) drop(s signal, n int) {
	p.meters.dropped.add(int64(n), slog.String("signal", signals[s]))
}

func (p *Pipeline) run() {
	defer close(p.loop)
	flush := time.NewTicker(flushEvery)
	collect := time.NewTicker(collectEvery)
	defer flush.Stop()
	defer collect.Stop()
	for {
		select {
		case <-p.stop:
			return
		case <-p.kick:
			p.flush(time.Now(), false)
		case <-flush.C:
			p.flush(time.Now(), false)
		case <-collect.C:
			p.collect(time.Now(), false)
		}
	}
}

func (p *Pipeline) flush(now time.Time, final bool) {
	if p.ready(sigTraces, now, final) {
		if spans := p.spans.take(); len(spans) > 0 {
			p.deliver(sigTraces, encodeTraces(nil, p.res, spans), len(spans), now, final)
		}
	}
	if p.ready(sigLogs, now, final) {
		if recs := p.recs.take(); len(recs) > 0 {
			p.deliver(sigLogs, encodeLogs(nil, p.res, recs), len(recs), now, final)
		}
	}
	p.ready(sigMetrics, now, final)
}

func (p *Pipeline) collect(now time.Time, final bool) {
	if !p.ready(sigMetrics, now, final) {
		return
	}
	ms := p.meters.gather()
	n := 0
	for _, m := range ms {
		n += len(m.points)
	}
	if n > 0 {
		p.deliver(sigMetrics, encodeMetrics(nil, p.res, p.start, now.UnixNano(), ms), n, now, final)
	}
}

func (p *Pipeline) ready(s signal, now time.Time, final bool) bool {
	r := p.retry[s]
	if r == nil {
		return true
	}
	if now.Before(r.due) {
		if !final {
			return false
		}
		p.retry[s] = nil
		p.drop(s, r.n)
		return true
	}
	p.retry[s] = nil
	if !p.send(s, r.body).ok {
		p.drop(s, r.n)
	}
	return true
}

func (p *Pipeline) deliver(s signal, body []byte, n int, now time.Time, final bool) {
	var buf bytes.Buffer
	zw, _ := gzip.NewWriterLevel(&buf, gzip.BestSpeed)
	zw.Write(body)
	zw.Close()
	o := p.send(s, buf.Bytes())
	switch {
	case o.ok:
	case o.retry && !final && o.after <= maxRetry:
		p.retry[s] = &pending{body: buf.Bytes(), n: n, due: now.Add(max(o.after, minRetry))}
	default:
		p.drop(s, n)
	}
}

type outcome struct {
	ok, retry bool
	after     time.Duration
	why       string
}

func (p *Pipeline) send(s signal, body []byte) outcome {
	req, err := http.NewRequestWithContext(p.ctx, http.MethodPost, p.urls[s], bytes.NewReader(body))
	if err != nil {
		return p.note(s, outcome{why: err.Error()})
	}
	req.Header.Set("Content-Type", "application/x-protobuf")
	req.Header.Set("Content-Encoding", "gzip")
	resp, err := p.client.Do(req)
	if err != nil {
		return p.note(s, outcome{retry: p.ctx.Err() == nil, why: err.Error()})
	}
	io.Copy(io.Discard, io.LimitReader(resp.Body, 64<<10))
	resp.Body.Close()
	switch resp.StatusCode {
	case http.StatusTooManyRequests, http.StatusBadGateway, http.StatusServiceUnavailable, http.StatusGatewayTimeout:
		return p.note(s, outcome{retry: true, after: retryAfter(resp.Header.Get("Retry-After")), why: resp.Status})
	}
	if resp.StatusCode/100 != 2 {
		return p.note(s, outcome{why: resp.Status})
	}
	return p.note(s, outcome{ok: true})
}

func (p *Pipeline) note(s signal, o outcome) outcome {
	h := p.warn.Load()
	if h == nil || p.ctx.Err() != nil || o.ok != p.failing[s] {
		return o
	}
	p.failing[s] = !o.ok
	rec := slog.NewRecord(time.Now(), slog.LevelInfo, "telemetry export recovered", 0)
	if !o.ok {
		rec = slog.NewRecord(time.Now(), slog.LevelWarn, "telemetry export failing", 0)
		rec.AddAttrs(slog.String("reason", o.why))
	}
	rec.AddAttrs(slog.String("signal", signals[s]), slog.String("endpoint", p.paths[s]))
	if (*h).Enabled(context.Background(), rec.Level) {
		(*h).Handle(context.Background(), rec)
	}
	return o
}

func (p *Pipeline) local(h slog.Handler) slog.Handler {
	if l, ok := h.(*logs); ok {
		h = l.local
	}
	if p != nil && h != nil {
		p.warn.CompareAndSwap(nil, &h)
	}
	return h
}

func retryAfter(h string) time.Duration {
	if s, err := strconv.Atoi(h); err == nil {
		if s > int(maxRetry/time.Second) {
			return maxRetry + time.Second
		}
		return time.Duration(max(s, 0)) * time.Second
	}
	if t, err := http.ParseTime(h); err == nil {
		return max(time.Until(t), 0)
	}
	return 0
}

func (p *Pipeline) Shutdown(ctx context.Context) error {
	if p == nil {
		return nil
	}
	if p.closed.CompareAndSwap(false, true) {
		release := context.AfterFunc(ctx, p.cancel)
		close(p.stop)
		<-p.loop
		now := time.Now()
		p.flush(now, true)
		p.collect(now, true)
		release()
		p.cancel()
		p.err = ctx.Err()
		close(p.done)
		return p.err
	}
	select {
	case <-p.done:
		return p.err
	case <-ctx.Done():
		return ctx.Err()
	}
}
