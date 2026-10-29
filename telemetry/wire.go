package telemetry

import (
	"encoding/binary"
	"log/slog"
	"math"
	"strings"
	"time"
	"unicode/utf8"
)

type enc struct{ b []byte }

func (e *enc) tag(f, wt int) { e.b = binary.AppendUvarint(e.b, uint64(f)<<3|uint64(wt)) }

func (e *enc) varint(f int, v uint64) {
	if v == 0 {
		return
	}
	e.tag(f, 0)
	e.b = binary.AppendUvarint(e.b, v)
}

func (e *enc) fixed64(f int, v uint64) {
	if v == 0 {
		return
	}
	e.tag(f, 1)
	e.b = binary.LittleEndian.AppendUint64(e.b, v)
}

func (e *enc) fixed32(f int, v uint32) {
	if v == 0 {
		return
	}
	e.tag(f, 5)
	e.b = binary.LittleEndian.AppendUint32(e.b, v)
}

func (e *enc) double(f int, v float64) {
	e.tag(f, 1)
	e.b = binary.LittleEndian.AppendUint64(e.b, math.Float64bits(v))
}

func valid(s string) string {
	if utf8.ValidString(s) {
		return s
	}
	return strings.ToValidUTF8(s, "\uFFFD")
}

func (e *enc) str(f int, s string) {
	if s == "" {
		return
	}
	s = valid(s)
	e.tag(f, 2)
	e.b = binary.AppendUvarint(e.b, uint64(len(s)))
	e.b = append(e.b, s...)
}

func (e *enc) bytes(f int, p []byte) {
	if len(p) == 0 {
		return
	}
	e.tag(f, 2)
	e.b = binary.AppendUvarint(e.b, uint64(len(p)))
	e.b = append(e.b, p...)
}

func (e *enc) open(f int) int {
	e.tag(f, 2)
	e.b = append(e.b, 0)
	return len(e.b)
}

func (e *enc) close(start int) {
	n := len(e.b) - start
	if n < 0x80 {
		e.b[start-1] = byte(n)
		return
	}
	var tmp [binary.MaxVarintLen64]byte
	l := binary.PutUvarint(tmp[:], uint64(n))
	e.b = append(e.b, tmp[:l-1]...)
	copy(e.b[start-1+l:], e.b[start:start+n])
	copy(e.b[start-1:], tmp[:l])
}

func (e *enc) oneofString(f int, s string) {
	s = valid(s)
	e.tag(f, 2)
	e.b = binary.AppendUvarint(e.b, uint64(len(s)))
	e.b = append(e.b, s...)
}

func (e *enc) anyValue(f int, v slog.Value) {
	m := e.open(f)
	switch v.Kind() {
	case slog.KindString:
		e.oneofString(1, v.String())
	case slog.KindBool:
		e.tag(2, 0)
		var b uint64
		if v.Bool() {
			b = 1
		}
		e.b = binary.AppendUvarint(e.b, b)
	case slog.KindInt64:
		e.tag(3, 0)
		e.b = binary.AppendUvarint(e.b, uint64(v.Int64()))
	case slog.KindUint64:
		e.tag(3, 0)
		e.b = binary.AppendUvarint(e.b, v.Uint64())
	case slog.KindFloat64:
		e.double(4, v.Float64())
	case slog.KindDuration:
		e.tag(3, 0)
		e.b = binary.AppendUvarint(e.b, uint64(v.Duration().Nanoseconds()))
	case slog.KindTime:
		e.oneofString(1, v.Time().Format(time.RFC3339Nano))
	case slog.KindGroup:
		kv := e.open(6)
		for _, a := range v.Group() {
			e.keyValue(1, a)
		}
		e.close(kv)
	default:
		e.oneofString(1, v.String())
	}
	e.close(m)
}

func (e *enc) keyValue(f int, a slog.Attr) {
	m := e.open(f)
	e.str(1, a.Key)
	e.anyValue(2, a.Value.Resolve())
	e.close(m)
}

func (e *enc) resource(res []byte) {
	e.bytes(1, res)
}

func (e *enc) scope(f int) {
	m := e.open(f)
	e.str(1, scope)
	e.str(2, version)
	e.close(m)
}

func encodeResource(attrs []slog.Attr) []byte {
	var e enc
	for _, a := range attrs {
		e.keyValue(1, a)
	}
	return e.b
}

const sampled = 0x01

func encodeTraces(dst, res []byte, spans []*Span) []byte {
	e := enc{b: dst}
	rs := e.open(1)
	e.resource(res)
	ss := e.open(2)
	e.scope(1)
	for _, s := range spans {
		sp := e.open(2)
		e.bytes(1, s.trace[:])
		e.bytes(2, s.id[:])
		if s.parent != (spanID{}) {
			e.bytes(4, s.parent[:])
		}
		e.str(5, s.name)
		e.varint(6, uint64(s.kind))
		e.fixed64(7, uint64(s.start))
		e.fixed64(8, uint64(s.end))
		for _, a := range s.attrs {
			e.keyValue(9, a)
		}
		e.varint(10, uint64(s.dropped))
		for _, l := range s.links {
			lk := e.open(13)
			e.bytes(1, l.trace[:])
			e.bytes(2, l.id[:])
			e.fixed32(6, 0x300|uint32(l.flags&sampled))
			e.close(lk)
		}
		if s.code != 0 {
			st := e.open(15)
			e.str(2, s.msg)
			e.varint(3, uint64(s.code))
			e.close(st)
		}
		e.fixed32(16, 0x100|sampled)
		e.close(sp)
	}
	e.close(ss)
	e.close(rs)
	return e.b
}

func encodeLogs(dst, res []byte, logs []record) []byte {
	e := enc{b: dst}
	rl := e.open(1)
	e.resource(res)
	sl := e.open(2)
	e.scope(1)
	for i := range logs {
		r := &logs[i]
		lr := e.open(2)
		e.fixed64(1, uint64(r.time))
		e.varint(2, uint64(severity(r.level)))
		e.str(3, r.level.String())
		e.anyValue(5, slog.StringValue(r.msg))
		for _, a := range r.attrs {
			e.keyValue(6, a)
		}
		if r.trace != (traceID{}) {
			e.fixed32(8, sampled)
			e.bytes(9, r.trace[:])
			e.bytes(10, r.span[:])
		}
		e.fixed64(11, uint64(r.observed))
		e.str(12, r.event)
		e.close(lr)
	}
	e.close(sl)
	e.close(rl)
	return e.b
}

const cumulative = 2

func encodeMetrics(dst, res []byte, start, now int64, ms []metric) []byte {
	e := enc{b: dst}
	rm := e.open(1)
	e.resource(res)
	sm := e.open(2)
	e.scope(1)
	for i := range ms {
		m := &ms[i]
		mm := e.open(2)
		e.str(1, m.name)
		e.str(3, m.unit)
		if m.kind == histogram {
			h := e.open(9)
			for _, p := range m.points {
				dp := e.open(1)
				e.fixed64(2, uint64(start))
				e.fixed64(3, uint64(now))
				e.tag(4, 1)
				e.b = binary.LittleEndian.AppendUint64(e.b, p.count)
				e.double(5, p.sum)
				e.tag(6, 2)
				e.b = binary.AppendUvarint(e.b, uint64(8*len(p.buckets)))
				for _, c := range p.buckets {
					e.b = binary.LittleEndian.AppendUint64(e.b, c)
				}
				e.tag(7, 2)
				e.b = binary.AppendUvarint(e.b, uint64(8*len(m.bounds)))
				for _, x := range m.bounds {
					e.b = binary.LittleEndian.AppendUint64(e.b, math.Float64bits(x))
				}
				for _, a := range p.attrs {
					e.keyValue(9, a)
				}
				e.close(dp)
			}
			e.varint(2, cumulative)
			e.close(h)
		} else {
			s := e.open(7)
			for _, p := range m.points {
				dp := e.open(1)
				e.fixed64(2, uint64(start))
				e.fixed64(3, uint64(now))
				e.tag(6, 1)
				e.b = binary.LittleEndian.AppendUint64(e.b, uint64(p.value))
				for _, a := range p.attrs {
					e.keyValue(7, a)
				}
				e.close(dp)
			}
			e.varint(2, cumulative)
			if m.kind == counter {
				e.varint(3, 1)
			}
			e.close(s)
		}
		e.close(mm)
	}
	e.close(sm)
	e.close(rm)
	return e.b
}
