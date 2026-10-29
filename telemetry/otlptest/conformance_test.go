package otlptest

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"maps"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/cogikyo/hgmx/telemetry"
	commonpb "go.opentelemetry.io/proto/otlp/common/v1"
	logspb "go.opentelemetry.io/proto/otlp/logs/v1"
	metricspb "go.opentelemetry.io/proto/otlp/metrics/v1"
	tracepb "go.opentelemetry.io/proto/otlp/trace/v1"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protoreflect"
)

var bounds = []float64{0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10}

type sink struct {
	t       *testing.T
	mu      sync.Mutex
	traces  []*tracepb.TracesData
	logs    []*logspb.LogsData
	metrics []*metricspb.MetricsData
}

func (k *sink) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost || r.Header.Get("Content-Type") != "application/x-protobuf" || r.Header.Get("Content-Encoding") != "gzip" || r.URL.RawQuery != "tenant=x" {
		k.t.Errorf("%s %s: content type %q, encoding %q", r.Method, r.URL, r.Header.Get("Content-Type"), r.Header.Get("Content-Encoding"))
	}
	zr, err := gzip.NewReader(r.Body)
	if err != nil {
		k.t.Error(err)
		return
	}
	body, _ := io.ReadAll(zr)
	var m proto.Message
	k.mu.Lock()
	defer k.mu.Unlock()
	switch r.URL.Path {
	case "/otlp/v1/traces":
		x := &tracepb.TracesData{}
		k.traces, m = append(k.traces, x), x
	case "/otlp/v1/logs":
		x := &logspb.LogsData{}
		k.logs, m = append(k.logs, x), x
	case "/otlp/v1/metrics":
		x := &metricspb.MetricsData{}
		k.metrics, m = append(k.metrics, x), x
	default:
		k.t.Errorf("path %s", r.URL.Path)
		return
	}
	if err := proto.Unmarshal(body, m); err != nil {
		k.t.Errorf("%s: %v", r.URL.Path, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if again, _ := proto.Marshal(m); len(again) != len(body) {
		k.t.Errorf("%s: official re-marshal is %d bytes, received %d", r.URL.Path, len(again), len(body))
	}
	if unknown(m.ProtoReflect()) {
		k.t.Errorf("%s: unknown fields", r.URL.Path)
	}
}

func unknown(m protoreflect.Message) bool {
	if len(m.GetUnknown()) > 0 {
		return true
	}
	found := false
	m.Range(func(fd protoreflect.FieldDescriptor, v protoreflect.Value) bool {
		switch {
		case fd.IsList() && fd.Message() != nil:
			for i := range v.List().Len() {
				found = found || unknown(v.List().Get(i).Message())
			}
		case fd.Message() != nil && !fd.IsMap():
			found = found || unknown(v.Message())
		}
		return !found
	})
	return found
}

func flat(kvs []*commonpb.KeyValue) map[string]any {
	out := map[string]any{}
	for _, kv := range kvs {
		switch v := kv.Value.Value.(type) {
		case *commonpb.AnyValue_StringValue:
			out[kv.Key] = v.StringValue
		case *commonpb.AnyValue_IntValue:
			out[kv.Key] = v.IntValue
		case *commonpb.AnyValue_BoolValue:
			out[kv.Key] = v.BoolValue
		case *commonpb.AnyValue_DoubleValue:
			out[kv.Key] = v.DoubleValue
		case *commonpb.AnyValue_KvlistValue:
			out[kv.Key] = flat(v.KvlistValue.Values)
		default:
			out[kv.Key] = v
		}
	}
	return out
}

func same(t *testing.T, what string, got, want map[string]any) {
	t.Helper()
	if len(got) != len(want) {
		t.Errorf("%s attributes: got %v, want %v", what, got, want)
		return
	}
	for k, w := range want {
		g := got[k]
		if w == nil {
			if _, ok := got[k]; !ok {
				t.Errorf("%s: missing %s", what, k)
			}
			continue
		}
		if gm, ok := g.(map[string]any); ok {
			same(t, what+"."+k, gm, w.(map[string]any))
			continue
		}
		if g != w {
			t.Errorf("%s: %s = %#v, want %#v", what, k, g, w)
		}
	}
}

func within(t *testing.T, what string, ns uint64, lo, hi time.Time) {
	t.Helper()
	if ts := time.Unix(0, int64(ns)); ts.Before(lo) || ts.After(hi) {
		t.Errorf("%s %v is outside [%v, %v]", what, ts, lo, hi)
	}
}

func TestConformance(t *testing.T) {
	k := &sink{t: t}
	receiver := httptest.NewServer(k)
	defer receiver.Close()
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("traceparent") != "" || r.Header.Get("tracestate") != "" {
			t.Error("Transport forwarded trace headers to a third party")
		}
	}))
	defer upstream.Close()

	t.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", receiver.URL+"/otlp?tenant=x")
	t.Setenv("OTEL_EXPORTER_OTLP_PROTOCOL", "http/protobuf")
	t.Setenv("OTEL_SERVICE_NAME", "hgmx.dev")
	t.Setenv("OTEL_RESOURCE_ATTRIBUTES", "service.version=8ed2b22,host.name=apps,deployment.environment=production")
	lo := time.Now()
	p, err := telemetry.FromEnv()
	if err != nil || p == nil {
		t.Fatal(p, err)
	}

	var journal bytes.Buffer
	local := slog.NewTextHandler(&journal, nil)
	logger := slog.New(p.Logs(local))
	client := &http.Client{Transport: p.Transport(nil)}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /icons/{name}", func(w http.ResponseWriter, r *http.Request) {
		ctx, s := telemetry.Start(r.Context(), "render",
			slog.String("icon", r.PathValue("name")), slog.Int("n", 3), slog.Bool("cached", true),
			slog.Float64("ratio", 1.5), slog.Group("g", slog.String("k", "v")), slog.Int("n", 4))
		logger.InfoContext(ctx, "rendered", slog.Int("icons", 19))
		req, _ := http.NewRequestWithContext(ctx, "GET", upstream.URL+"/repos/cogikyo/hgmx?token=secret", nil)
		req.Header.Set("traceparent", "00-11111111111111111111111111111111-2222222222222222-01")
		if resp, err := client.Do(req); err == nil {
			resp.Body.Close()
		}
		time.Sleep(time.Millisecond)
		s.End()
		io.WriteString(w, "ok")
	})
	mux.HandleFunc("GET /boom", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	})
	beacon, err := p.Beacon()
	if err != nil {
		t.Fatal(err)
	}
	root := http.NewServeMux()
	root.Handle(telemetry.Path, beacon)
	root.Handle("/", p.Handler(local, mux))

	inbound := "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01"
	r := httptest.NewRequest("GET", "/icons/user", nil)
	r.Header.Set("traceparent", inbound)
	w := httptest.NewRecorder()
	root.ServeHTTP(w, r)
	timing := w.Header().Get("Server-Timing")
	tp := strings.TrimSuffix(strings.TrimPrefix(timing, `traceparent;desc="`), `"`)
	root.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/boom", nil))
	root.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/%FF", nil))

	ua := "Mozilla/5.0 (X11; Linux x86_64) Chrome/140.0 Safari/537.36"
	for _, body := range []string{
		`{"type":"view","path":"/icons/user","referrer":"https://news.ycombinator.com/item?id=1","traceparent":"` + tp + `"}`,
		`{"type":"vitals","path":"/icons/user","traceparent":"` + tp + `","ttfb":2,"fcp":28,"lcp":28,"cls":0.1694,"inp":120}`,
	} {
		req := httptest.NewRequest("POST", telemetry.Path, strings.NewReader(body))
		req.Header.Set("User-Agent", ua)
		req.Header.Set("CF-IPCountry", "NZ")
		bw := httptest.NewRecorder()
		root.ServeHTTP(bw, req)
		if bw.Code != http.StatusNoContent {
			t.Errorf("beacon: %d", bw.Code)
		}
	}
	if resp, err := client.Get(upstream.URL + "/background"); err == nil {
		resp.Body.Close()
	}
	_, job := p.Start(context.Background(), "job")
	job.Fail(errors.New("nightly failed"))
	job.End()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := p.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}
	hi := time.Now()

	if len(k.traces) != 1 || len(k.logs) != 1 || len(k.metrics) != 1 {
		t.Fatalf("payloads: %d traces, %d logs, %d metrics; want 1 each", len(k.traces), len(k.logs), len(k.metrics))
	}
	resource := map[string]any{
		"service.name": "hgmx.dev", "service.version": "8ed2b22", "host.name": "apps", "deployment.environment": "production",
		"telemetry.sdk.name": "hgmx", "telemetry.sdk.language": "go", "telemetry.sdk.version": nil,
	}
	rs := k.traces[0].ResourceSpans
	rl := k.logs[0].ResourceLogs
	rm := k.metrics[0].ResourceMetrics
	if len(rs) != 1 || len(rl) != 1 || len(rm) != 1 || len(rs[0].ScopeSpans) != 1 || len(rl[0].ScopeLogs) != 1 || len(rm[0].ScopeMetrics) != 1 {
		t.Fatal("want one resource and one scope per payload")
	}
	for name, res := range map[string]*struct {
		attrs  []*commonpb.KeyValue
		scope  *commonpb.InstrumentationScope
		schema string
	}{
		"traces":  {rs[0].Resource.Attributes, rs[0].ScopeSpans[0].Scope, rs[0].SchemaUrl + rs[0].ScopeSpans[0].SchemaUrl},
		"logs":    {rl[0].Resource.Attributes, rl[0].ScopeLogs[0].Scope, rl[0].SchemaUrl + rl[0].ScopeLogs[0].SchemaUrl},
		"metrics": {rm[0].Resource.Attributes, rm[0].ScopeMetrics[0].Scope, rm[0].SchemaUrl + rm[0].ScopeMetrics[0].SchemaUrl},
	} {
		same(t, name+" resource", flat(res.attrs), resource)
		if res.scope.Name != "github.com/cogikyo/hgmx/telemetry" || len(res.scope.Attributes) != 0 || res.schema != "" {
			t.Errorf("%s scope %v", name, res.scope)
		}
	}

	spans := map[string]*tracepb.Span{}
	for _, s := range rs[0].ScopeSpans[0].Spans {
		name := s.Name
		if s.Kind == tracepb.Span_SPAN_KIND_CLIENT && len(s.ParentSpanId) == 0 {
			name = "background"
		}
		if s.Name == "GET" && s.Kind == tracepb.Span_SPAN_KIND_SERVER {
			name = "unrouted"
		}
		spans[name] = s
		within(t, name+" start", s.StartTimeUnixNano, lo, hi)
		within(t, name+" end", s.EndTimeUnixNano, lo, hi)
		if s.StartTimeUnixNano > s.EndTimeUnixNano || len(s.TraceId) != 16 || len(s.SpanId) != 8 || s.TraceState != "" ||
			s.Flags != 0x101 || s.DroppedAttributesCount != 0 || len(s.Events) != 0 || s.DroppedEventsCount != 0 || s.DroppedLinksCount != 0 {
			t.Errorf("span %s: %v", name, s)
		}
	}
	if len(spans) != 7 {
		t.Fatalf("spans: %v", slices.Collect(maps.Keys(spans)))
	}
	srv, render, call, bg, boom, unrouted, job2 := spans["GET /icons/{name}"], spans["render"], spans["GET"], spans["background"], spans["GET /boom"], spans["unrouted"], spans["job"]
	if srv.EndTimeUnixNano-srv.StartTimeUnixNano < uint64(time.Millisecond) {
		t.Errorf("server span lasted %dns; the handler sleeps 1ms", srv.EndTimeUnixNano-srv.StartTimeUnixNano)
	}

	if srv.Kind != tracepb.Span_SPAN_KIND_SERVER || len(srv.ParentSpanId) != 0 || srv.Status != nil {
		t.Errorf("server span: %v", srv)
	}
	same(t, "server span", flat(srv.Attributes), map[string]any{
		"http.request.method": "GET", "http.route": "/icons/{name}", "http.response.status_code": int64(200),
		"url.path": "/icons/user", "http.response.body.size": int64(2),
	})
	if len(srv.Links) != 1 || hex.EncodeToString(srv.Links[0].TraceId) != inbound[3:35] || hex.EncodeToString(srv.Links[0].SpanId) != inbound[36:52] ||
		srv.Links[0].Flags != 0x301 || srv.Links[0].TraceState != "" || len(srv.Links[0].Attributes) != 0 {
		t.Errorf("inbound traceparent link: %v", srv.Links)
	}
	if timing != `traceparent;desc="00-`+hex.EncodeToString(srv.TraceId)+"-"+hex.EncodeToString(srv.SpanId)+`-01"` {
		t.Errorf("Server-Timing %q", timing)
	}

	if render.Kind != tracepb.Span_SPAN_KIND_INTERNAL || !bytes.Equal(render.ParentSpanId, srv.SpanId) || !bytes.Equal(render.TraceId, srv.TraceId) ||
		render.StartTimeUnixNano < srv.StartTimeUnixNano || render.EndTimeUnixNano > srv.EndTimeUnixNano || render.Status != nil || len(render.Links) != 0 {
		t.Errorf("render span: %v", render)
	}
	same(t, "render span", flat(render.Attributes), map[string]any{
		"icon": "user", "n": int64(4), "cached": true, "ratio": 1.5, "g": map[string]any{"k": "v"},
	})

	if call.Kind != tracepb.Span_SPAN_KIND_CLIENT || !bytes.Equal(call.ParentSpanId, render.SpanId) || !bytes.Equal(call.TraceId, srv.TraceId) ||
		call.StartTimeUnixNano < render.StartTimeUnixNano || call.EndTimeUnixNano > render.EndTimeUnixNano || call.Status != nil {
		t.Errorf("client span: %v", call)
	}
	host, port, _ := strings.Cut(strings.TrimPrefix(upstream.URL, "http://"), ":")
	same(t, "client span", flat(call.Attributes), map[string]any{
		"http.request.method": "GET", "server.address": host, "server.port": nil,
		"url.full": "http://" + host + ":" + port + "/repos/cogikyo/hgmx", "http.response.status_code": int64(200),
	})
	if bg.Kind != tracepb.Span_SPAN_KIND_CLIENT || bytes.Equal(bg.TraceId, srv.TraceId) || flat(bg.Attributes)["url.full"] != "http://"+host+":"+port+"/background" {
		t.Errorf("background client span: %v", bg)
	}

	if boom.Kind != tracepb.Span_SPAN_KIND_SERVER || boom.Status.GetCode() != tracepb.Status_STATUS_CODE_ERROR || boom.Status.Message != "" {
		t.Errorf("500 span: %v", boom)
	}
	same(t, "500 span", flat(boom.Attributes), map[string]any{
		"http.request.method": "GET", "http.route": "/boom", "http.response.status_code": int64(500),
		"url.path": "/boom", "http.response.body.size": int64(5),
	})
	if unrouted.Status != nil {
		t.Errorf("404 span status: %v; 4xx leaves server spans unset", unrouted.Status)
	}
	same(t, "unrouted span", flat(unrouted.Attributes), map[string]any{
		"http.request.method": "GET", "http.response.status_code": int64(404),
		"url.path": "/\uFFFD", "http.response.body.size": int64(19),
	})
	if job2.Kind != tracepb.Span_SPAN_KIND_INTERNAL || len(job2.ParentSpanId) != 0 || job2.Status.GetCode() != tracepb.Status_STATUS_CODE_ERROR || job2.Status.Message != "nightly failed" {
		t.Errorf("background job span: %v", job2)
	}

	records := map[string]*logspb.LogRecord{}
	for _, l := range rl[0].ScopeLogs[0].LogRecords {
		name := l.Body.GetStringValue()
		records[name] = l
		within(t, name+" time", l.TimeUnixNano, lo, hi)
		within(t, name+" observed", l.ObservedTimeUnixNano, lo, hi)
		event := ""
		if strings.HasPrefix(name, "browser.") {
			event = name
		}
		if l.ObservedTimeUnixNano < l.TimeUnixNano || l.SeverityNumber != logspb.SeverityNumber_SEVERITY_NUMBER_INFO || l.SeverityText != "INFO" ||
			l.Flags != 1 || l.DroppedAttributesCount != 0 || l.EventName != event {
			t.Errorf("log %s: %v", name, l)
		}
	}
	if len(records) != 3 {
		t.Fatalf("log records: %v", records)
	}
	rendered, view, vitals := records["rendered"], records["browser.page_view"], records["browser.vitals"]
	if !bytes.Equal(rendered.TraceId, render.TraceId) || !bytes.Equal(rendered.SpanId, render.SpanId) {
		t.Errorf("rendered ids: %v", rendered)
	}
	same(t, "rendered", flat(rendered.Attributes), map[string]any{"icons": int64(19)})
	for _, l := range []*logspb.LogRecord{view, vitals} {
		if !bytes.Equal(l.TraceId, srv.TraceId) || !bytes.Equal(l.SpanId, srv.SpanId) || l.TimeUnixNano != l.ObservedTimeUnixNano {
			t.Errorf("%s ids/time: %v", l.EventName, l)
		}
	}
	same(t, "page view", flat(view.Attributes), map[string]any{
		"url.path": "/icons/user", "referrer.host": "news.ycombinator.com", "visitor.hash": nil,
		"user_agent.name": "Chrome", "browser.mobile": false, "geo.country.iso_code": "NZ",
	})
	same(t, "vitals", flat(vitals.Attributes), map[string]any{
		"url.path": "/icons/user", "browser.vitals.ttfb": 2.0, "browser.vitals.fcp": 28.0, "browser.vitals.lcp": 28.0,
		"browser.vitals.inp": 120.0, "browser.vitals.cls": 0.1694,
		"user_agent.name": "Chrome", "browser.mobile": false, "geo.country.iso_code": "NZ",
	})
	if n := strings.Count(journal.String(), "trace_id="+hex.EncodeToString(srv.TraceId)); n != 2 {
		t.Errorf("journal lines with the trace id: %d, want 2 (access + rendered)\n%s", n, journal.String())
	}

	metrics := map[string]*metricspb.Metric{}
	for _, m := range rm[0].ScopeMetrics[0].Metrics {
		metrics[m.Name] = m
		if m.Description != "" || len(m.Metadata) != 0 {
			t.Errorf("metric %s: %v", m.Name, m)
		}
	}
	var start uint64
	d := metrics["http.server.request.duration"]
	h := d.GetHistogram()
	if d.Unit != "s" || h == nil || h.AggregationTemporality != metricspb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE || len(h.DataPoints) != 3 {
		t.Fatalf("duration: %v", d)
	}
	counts := map[string]uint64{}
	for _, dp := range h.DataPoints {
		a := flat(dp.Attributes)
		route, _ := a["http.route"].(string)
		counts[fmt.Sprint(route, " ", a["http.response.status_code"])] = dp.Count
		var n uint64
		for _, c := range dp.BucketCounts {
			n += c
		}
		start = dp.StartTimeUnixNano
		within(t, "duration start", dp.StartTimeUnixNano, lo, hi)
		within(t, "duration time", dp.TimeUnixNano, lo, hi)
		if !slices.Equal(dp.ExplicitBounds, bounds) || len(dp.BucketCounts) != 17 || n != dp.Count || dp.GetSum() <= 0 ||
			dp.StartTimeUnixNano >= dp.TimeUnixNano || dp.Flags != 0 || len(dp.Exemplars) != 0 || dp.Min != nil || dp.Max != nil || a["http.request.method"] != "GET" {
			t.Errorf("duration point: %v", dp)
		}
	}
	if counts["/icons/{name} 200"] != 1 || counts["/boom 500"] != 1 || counts[" 404"] != 1 {
		t.Errorf("duration counts: %v", counts)
	}
	for name, want := range map[string]struct {
		unit      string
		monotonic bool
		zero      bool
	}{
		"http.server.active_requests": {"{request}", false, true},
		"go.memory.used":              {"By", false, false},
		"go.memory.allocated":         {"By", true, false},
		"go.memory.gc.goal":           {"By", false, false},
		"go.goroutine.count":          {"{goroutine}", false, false},
	} {
		m := metrics[name]
		s := m.GetSum()
		if s == nil || m.Unit != want.unit || s.IsMonotonic != want.monotonic || len(s.DataPoints) != 1 ||
			s.AggregationTemporality != metricspb.AggregationTemporality_AGGREGATION_TEMPORALITY_CUMULATIVE {
			t.Errorf("%s: %v", name, m)
			continue
		}
		dp := s.DataPoints[0]
		if (dp.GetAsInt() == 0) != want.zero || len(dp.Attributes) != 0 || dp.StartTimeUnixNano != start || dp.TimeUnixNano <= start {
			t.Errorf("%s point: %v", name, dp)
		}
	}
	if len(metrics) != 6 {
		t.Errorf("metrics: %d, want 6 (no drops)", len(metrics))
	}

	if dir := os.Getenv("OTLP_FIXTURES"); dir != "" {
		opts := protojson.MarshalOptions{Multiline: true, Indent: "  "}
		for name, m := range map[string]proto.Message{"traces.json": k.traces[0], "logs.json": k.logs[0], "metrics.json": k.metrics[0]} {
			os.WriteFile(filepath.Join(dir, name), []byte(opts.Format(m)), 0o644)
		}
	}
}
