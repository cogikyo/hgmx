package telemetry

import (
	"context"
	"encoding/binary"
	"encoding/hex"
	"log/slog"
	"math/rand/v2"
	"slices"
	"sync"
	"time"
)

type (
	traceID [16]byte
	spanID  [8]byte
)

type spanKind uint8

const (
	internal spanKind = 1
	server   spanKind = 2
	client   spanKind = 3
)

const maxAttrs = 128

type link struct {
	trace traceID
	id    spanID
	flags byte
}

type Span struct {
	p       *Pipeline
	trace   traceID
	id      spanID
	parent  spanID
	name    string
	kind    spanKind
	start   int64
	end     int64
	attrs   []slog.Attr
	dropped int
	links   []link
	code    uint8
	msg     string
}

type spanKey struct{}

func SpanFrom(ctx context.Context) *Span {
	s, _ := ctx.Value(spanKey{}).(*Span)
	return s
}

func (s *Span) recording() bool { return s != nil && s.p != nil }

func Start(ctx context.Context, name string, attrs ...slog.Attr) (context.Context, *Span) {
	parent := SpanFrom(ctx)
	if !parent.recording() {
		return ctx, nil
	}
	return parent.p.Start(ctx, name, attrs...)
}

func (p *Pipeline) Start(ctx context.Context, name string, attrs ...slog.Attr) (context.Context, *Span) {
	s := p.begin(ctx, name, internal)
	if !s.recording() {
		return ctx, nil
	}
	s.Set(attrs...)
	return context.WithValue(ctx, spanKey{}, s), s
}

func (p *Pipeline) begin(ctx context.Context, name string, kind spanKind) *Span {
	if p.off() {
		return nil
	}
	parent := SpanFrom(ctx)
	if parent == nil {
		return p.root(name, kind)
	}
	if !parent.recording() {
		return nil
	}
	s := &Span{p: p, trace: parent.trace, parent: parent.id, name: clip(name), kind: kind, start: time.Now().UnixNano()}
	binary.LittleEndian.PutUint64(s.id[:], rand.Uint64())
	return s
}

func (p *Pipeline) root(name string, kind spanKind) *Span {
	if p.off() {
		return nil
	}
	if !p.roots.take() {
		return &Span{}
	}
	s := &Span{p: p, name: clip(name), kind: kind, start: time.Now().UnixNano()}
	binary.LittleEndian.PutUint64(s.trace[:8], rand.Uint64())
	binary.LittleEndian.PutUint64(s.trace[8:], rand.Uint64())
	binary.LittleEndian.PutUint64(s.id[:], rand.Uint64())
	return s
}

func (s *Span) Set(attrs ...slog.Attr) {
	if !s.recording() || s.end != 0 {
		return
	}
	flat := make([]slog.Attr, 0, len(attrs))
	for _, a := range attrs {
		flat = snapshot(flat, a)
	}
	s.attrs = slices.Grow(s.attrs, min(len(flat), maxAttrs-len(s.attrs)))
	for _, a := range flat {
		if i := s.index(a.Key); i >= 0 {
			s.attrs[i] = a
		} else if len(s.attrs) < maxAttrs {
			s.attrs = append(s.attrs, a)
		} else {
			s.dropped++
		}
	}
}

func (s *Span) index(key string) int {
	for i, a := range s.attrs {
		if a.Key == key {
			return i
		}
	}
	return -1
}

func (s *Span) Fail(err error) {
	if !s.recording() || s.end != 0 || err == nil {
		return
	}
	s.code, s.msg = 2, clip(err.Error())
}

func (s *Span) End() {
	if !s.recording() || s.end != 0 {
		return
	}
	s.end = time.Now().UnixNano()
	s.p.spans.add(s)
}

func (s *Span) TraceID() string {
	if !s.recording() {
		return ""
	}
	return hex.EncodeToString(s.trace[:])
}

func (s *Span) spanHex() string {
	return hex.EncodeToString(s.id[:])
}

func (s *Span) traceparent() string {
	var b [55]byte
	copy(b[:], "00-")
	hex.Encode(b[3:35], s.trace[:])
	b[35] = '-'
	hex.Encode(b[36:52], s.id[:])
	copy(b[52:], "-01")
	return string(b[:])
}

func lowerHex(s string) bool {
	for i := range len(s) {
		if c := s[i]; (c < '0' || c > '9') && (c < 'a' || c > 'f') {
			return false
		}
	}
	return true
}

func parseTraceparent(h string) (link, bool) {
	var l link
	if len(h) < 55 || h[2] != '-' || h[35] != '-' || h[52] != '-' || (len(h) > 55 && h[55] != '-') {
		return l, false
	}
	ver, trace, id, flags := h[:2], h[3:35], h[36:52], h[53:55]
	if !lowerHex(ver) || !lowerHex(trace) || !lowerHex(id) || !lowerHex(flags) || ver == "ff" || (ver == "00" && len(h) != 55) {
		return l, false
	}
	hex.Decode(l.trace[:], []byte(trace))
	hex.Decode(l.id[:], []byte(id))
	var f [1]byte
	hex.Decode(f[:], []byte(flags))
	l.flags = f[0]
	if l.trace == (traceID{}) || l.id == (spanID{}) {
		return l, false
	}
	return l, true
}

type bucket struct {
	mu     sync.Mutex
	rate   float64
	tokens float64
	last   time.Time
}

func newBucket(rate float64) *bucket {
	return &bucket{rate: rate, tokens: rate, last: time.Now()}
}

func (b *bucket) take() bool {
	b.mu.Lock()
	defer b.mu.Unlock()
	now := time.Now()
	b.tokens = min(b.rate, b.tokens+now.Sub(b.last).Seconds()*b.rate)
	b.last = now
	if b.tokens < 1 {
		return false
	}
	b.tokens--
	return true
}
