package telemetry

import (
	"context"
	"log/slog"
	"slices"
	"time"
)

type record struct {
	time     int64
	observed int64
	level    slog.Level
	msg      string
	event    string
	attrs    []slog.Attr
	trace    traceID
	span     spanID
}

func severity(l slog.Level) int {
	return min(max(int(l)+9, 1), 24)
}

type op struct {
	group string
	attrs []slog.Attr
}

type logs struct {
	local  slog.Handler
	base   slog.Handler
	ops    []op
	p      *Pipeline
	groups []string
	frames [][]slog.Attr
}

func (p *Pipeline) Logs(local slog.Handler) slog.Handler {
	if p == nil {
		return local
	}
	return &logs{local: local, base: local, p: p, frames: [][]slog.Attr{nil}}
}

func (h *logs) Enabled(ctx context.Context, l slog.Level) bool {
	return (l >= slog.LevelInfo && !h.p.off()) || h.local.Enabled(ctx, l)
}

func (h *logs) Handle(ctx context.Context, r slog.Record) error {
	s := SpanFrom(ctx)
	if r.Level >= slog.LevelInfo && !h.p.off() {
		h.p.recs.add(h.record(r, s))
	}
	if !h.local.Enabled(ctx, r.Level) {
		return nil
	}
	if !s.recording() {
		return h.local.Handle(ctx, r)
	}
	ids := []slog.Attr{slog.String("trace_id", s.TraceID()), slog.String("span_id", s.spanHex())}
	if len(h.groups) == 0 {
		r = r.Clone()
		r.AddAttrs(ids...)
		return h.local.Handle(ctx, r)
	}
	l := h.base.WithAttrs(ids)
	for _, o := range h.ops {
		if o.group != "" {
			l = l.WithGroup(o.group)
		} else {
			l = l.WithAttrs(o.attrs)
		}
	}
	return l.Handle(ctx, r)
}

func (h *logs) record(r slog.Record, s *Span) record {
	inner := slices.Clone(h.frames[len(h.frames)-1])
	r.Attrs(func(a slog.Attr) bool {
		inner = snapshot(inner, a)
		return true
	})
	for i := len(h.groups) - 1; i >= 0; i-- {
		outer := slices.Clone(h.frames[i])
		if len(inner) > 0 {
			outer = append(outer, slog.Attr{Key: h.groups[i], Value: slog.GroupValue(inner...)})
		}
		inner = outer
	}
	rec := record{observed: time.Now().UnixNano(), level: r.Level, msg: clip(r.Message), attrs: inner}
	if !r.Time.IsZero() {
		rec.time = r.Time.UnixNano()
	}
	if s.recording() {
		rec.trace, rec.span = s.trace, s.id
	}
	return rec
}

func (h *logs) WithAttrs(as []slog.Attr) slog.Handler {
	if len(as) == 0 {
		return h
	}
	c := *h
	last := len(h.frames) - 1
	frame := slices.Clone(h.frames[last])
	for _, a := range as {
		frame = snapshot(frame, a)
	}
	c.frames = append(slices.Clip(h.frames[:last]), frame)
	if len(h.groups) > 0 {
		c.ops = append(slices.Clip(h.ops), op{attrs: slices.Clone(as)})
	}
	c.local = h.local.WithAttrs(as)
	if len(h.groups) == 0 {
		c.base = c.local
	}
	return &c
}

func (h *logs) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	c := *h
	c.groups = append(slices.Clip(h.groups), name)
	c.frames = append(slices.Clip(h.frames), nil)
	c.ops = append(slices.Clip(h.ops), op{group: name})
	c.local = h.local.WithGroup(name)
	return &c
}
