package telemetry

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"
)

type fakeRT struct {
	mu    sync.Mutex
	paths []string
	reply func(*http.Request) (*http.Response, error)
}

func (f *fakeRT) RoundTrip(r *http.Request) (*http.Response, error) {
	f.mu.Lock()
	f.paths = append(f.paths, r.URL.String())
	f.mu.Unlock()
	return f.reply(r)
}

func (f *fakeRT) calls() int {
	f.mu.Lock()
	defer f.mu.Unlock()
	return len(f.paths)
}

func status(code int, retryAfter string) func(*http.Request) (*http.Response, error) {
	return func(r *http.Request) (*http.Response, error) {
		h := http.Header{}
		if retryAfter != "" {
			h.Set("Retry-After", retryAfter)
		}
		if code/100 == 3 {
			h.Set("Location", "http://elsewhere.test/ingest")
		}
		return &http.Response{StatusCode: code, Status: strconv.Itoa(code) + " " + http.StatusText(code), Header: h, Body: io.NopCloser(http.NoBody), Request: r}, nil
	}
}

func offline(rt http.RoundTripper) *Pipeline {
	u, _ := url.Parse("http://receiver.test/otlp")
	p := newPipeline(u, nil, rt)
	p.roots = newBucket(1e9)
	return p
}

func spans(p *Pipeline, n int) {
	for range n {
		_, s := p.Start(context.Background(), "work")
		s.End()
	}
}

func dropped(p *Pipeline, sig signal) int64 {
	i := p.meters.dropped
	i.mu.Lock()
	defer i.mu.Unlock()
	for _, s := range i.series {
		if attr(s.attrs, "signal") == signals[sig] {
			return s.value
		}
	}
	return 0
}

func TestBufferDrops(t *testing.T) {
	p := offline(nil)
	spans(p, maxItems+10)
	if n := len(p.spans.take()); n != maxItems {
		t.Errorf("buffered %d spans, want %d", n, maxItems)
	}
	if d := dropped(p, sigTraces); d != 10 {
		t.Errorf("dropped %d spans, want 10", d)
	}
}

func TestShutdownAbortsInFlightExport(t *testing.T) {
	started := make(chan struct{})
	rt := &fakeRT{reply: func(r *http.Request) (*http.Response, error) {
		if r.URL.Path == "/otlp/v1/traces" {
			close(started)
		}
		<-r.Context().Done()
		return nil, r.Context().Err()
	}}
	p := offline(rt)
	go p.run()
	spans(p, batch)
	select {
	case <-started:
	case <-time.After(2 * time.Second):
		t.Fatal("export never started")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	begin := time.Now()
	err := p.Shutdown(ctx)
	took := time.Since(begin)
	if !errors.Is(err, context.DeadlineExceeded) || took > time.Second {
		t.Errorf("Shutdown with an export in flight: %v after %v", err, took)
	}
	if d := dropped(p, sigTraces); d != batch {
		t.Errorf("dropped %d spans, want the in-flight %d", d, batch)
	}
}

func TestRetryLifecycle(t *testing.T) {
	rt := &fakeRT{reply: status(http.StatusServiceUnavailable, "7")}
	p := offline(rt)
	t0 := time.Now()

	spans(p, 3)
	p.flush(t0, false)
	if rt.calls() != 1 || p.retry[sigTraces] == nil || !p.retry[sigTraces].due.Equal(t0.Add(7*time.Second)) {
		t.Fatalf("503 Retry-After 7: %d calls, retry %+v", rt.calls(), p.retry[sigTraces])
	}
	spans(p, 2)
	p.flush(t0.Add(6*time.Second), false)
	if rt.calls() != 1 {
		t.Errorf("sent %d times before the retry was due", rt.calls())
	}
	rt.reply = status(http.StatusOK, "")
	p.flush(t0.Add(7*time.Second), false)
	if rt.calls() != 3 || p.retry[sigTraces] != nil || dropped(p, sigTraces) != 0 {
		t.Errorf("at the retry time: %d calls, want the retry and the held batch", rt.calls())
	}

	rt.reply = status(http.StatusTooManyRequests, "1")
	spans(p, 1)
	p.flush(t0, false)
	if r := p.retry[sigTraces]; r == nil || !r.due.Equal(t0.Add(minRetry)) {
		t.Errorf("Retry-After 1 must wait %v: %+v", minRetry, r)
	}
	p.flush(t0.Add(time.Second), true)
	if rt.calls() != 4 || dropped(p, sigTraces) != 1 {
		t.Errorf("final flush before the retry time: %d calls, %d dropped; want no send, 1 dropped", rt.calls(), dropped(p, sigTraces))
	}

	rt.reply = status(http.StatusServiceUnavailable, "120")
	spans(p, 2)
	p.flush(t0, false)
	if p.retry[sigTraces] != nil || dropped(p, sigTraces) != 3 {
		t.Errorf("Retry-After 120: retry %+v, %d dropped; want dropped", p.retry[sigTraces], dropped(p, sigTraces))
	}

	rt.reply = status(http.StatusServiceUnavailable, "9999999999")
	spans(p, 1)
	p.flush(t0, false)
	if p.retry[sigTraces] != nil || dropped(p, sigTraces) != 4 {
		t.Errorf("Retry-After 9999999999: retry %+v; want dropped", p.retry[sigTraces])
	}

	for _, reply := range []func(*http.Request) (*http.Response, error){
		status(http.StatusBadGateway, ""),
		status(http.StatusGatewayTimeout, ""),
		func(*http.Request) (*http.Response, error) { return nil, errors.New("connection refused") },
	} {
		rt.reply = reply
		spans(p, 1)
		p.flush(t0, false)
		if r := p.retry[sigTraces]; r == nil || !r.due.Equal(t0.Add(minRetry)) {
			t.Errorf("502, 504 and transport errors must retry after %v: %+v", minRetry, r)
		}
		p.retry[sigTraces] = nil
	}

	rt.reply = status(http.StatusBadRequest, "")
	spans(p, 1)
	p.flush(t0, false)
	if p.retry[sigTraces] != nil || dropped(p, sigTraces) != 5 {
		t.Errorf("400 must drop without retry")
	}

	rt.reply = status(http.StatusServiceUnavailable, "")
	calls := rt.calls()
	p.collect(t0, false)
	p.collect(t0.Add(time.Second), false)
	if rt.calls() != calls+1 || p.retry[sigMetrics] == nil {
		t.Errorf("metrics 503: %d calls, retry %+v", rt.calls()-calls, p.retry[sigMetrics])
	}
	rt.reply = status(http.StatusOK, "")
	p.flush(t0.Add(minRetry), false)
	if rt.calls() != calls+2 || p.retry[sigMetrics] != nil {
		t.Errorf("metrics retry at its time: %d calls", rt.calls()-calls)
	}
}

func TestExportWarningsAreEdges(t *testing.T) {
	rt := &fakeRT{reply: status(http.StatusNotFound, "")}
	p := offline(rt)
	var out bytes.Buffer
	p.Handler(slog.NewTextHandler(&out, nil), http.NotFoundHandler())
	for range 3 {
		spans(p, 1)
		p.flush(time.Now(), false)
	}
	rt.reply = status(http.StatusOK, "")
	for range 2 {
		spans(p, 1)
		p.flush(time.Now(), false)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != 2 || !strings.Contains(lines[0], `msg="telemetry export failing" reason="404 Not Found" signal=spans endpoint=/otlp/v1/traces`) ||
		!strings.Contains(lines[1], `msg="telemetry export recovered" signal=spans endpoint=/otlp/v1/traces`) {
		t.Errorf("want one failing and one recovered line:\n%s", out.String())
	}
}

func TestRedirectNotFollowed(t *testing.T) {
	rt := &fakeRT{reply: status(http.StatusTemporaryRedirect, "")}
	p := offline(rt)
	spans(p, 1)
	p.flush(time.Now(), false)
	if rt.calls() != 1 || dropped(p, sigTraces) != 1 || p.retry[sigTraces] != nil {
		t.Errorf("307: %v, %d dropped", rt.paths, dropped(p, sigTraces))
	}
}

func TestShutdownOnce(t *testing.T) {
	rt := &fakeRT{reply: status(http.StatusOK, "")}
	p := offline(rt)
	go p.run()
	var wg sync.WaitGroup
	for range 3 {
		wg.Go(func() {
			if err := p.Shutdown(context.Background()); err != nil {
				t.Error(err)
			}
		})
	}
	wg.Wait()

	h := p.Handler(nil, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	h.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil))
	slog.New(p.Logs(slog.DiscardHandler)).Info("late")
	if _, s := p.Start(context.Background(), "late"); s != nil {
		t.Error("Start after Shutdown returned a span")
	}
	if n := len(p.spans.items) + len(p.recs.items) + len(p.meters.duration.series); n != 0 {
		t.Errorf("recorded %d items after Shutdown", n)
	}
}

func BenchmarkHandler(b *testing.B) {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /icons", func(w http.ResponseWriter, r *http.Request) { io.WriteString(w, "ok") })
	for _, c := range []struct {
		name string
		p    *Pipeline
	}{{"off", nil}, {"on", offline(&fakeRT{reply: status(http.StatusOK, "")})}} {
		b.Run(c.name, func(b *testing.B) {
			h := c.p.Handler(nil, mux)
			r := httptest.NewRequest("GET", "/icons", nil)
			b.ReportAllocs()
			for i := 0; b.Loop(); i++ {
				h.ServeHTTP(httptest.NewRecorder(), r)
				if c.p != nil && i%batch == 0 {
					c.p.spans.take()
				}
			}
		})
	}
}
