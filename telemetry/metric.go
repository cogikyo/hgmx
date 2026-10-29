package telemetry

import (
	"encoding/binary"
	"log/slog"
	"runtime/metrics"
	"slices"
	"strconv"
	"sync"
)

var bounds = []float64{
	0.0001, 0.00025, 0.0005, 0.001, 0.0025, 0.005, 0.01, 0.025,
	0.05, 0.1, 0.25, 0.5, 1, 2.5, 5, 10,
}

const maxSeries = 256

type kind uint8

const (
	counter kind = iota
	updown
	histogram
)

type series struct {
	attrs   []slog.Attr
	value   int64
	count   uint64
	sum     float64
	buckets []uint64
}

type instrument struct {
	name, unit string
	kind       kind
	bounds     []float64

	mu     sync.Mutex
	series map[string]*series
}

func newInstrument(name, unit string, k kind) *instrument {
	i := &instrument{name: name, unit: unit, kind: k, series: map[string]*series{}}
	if k == histogram {
		i.bounds = bounds
	}
	return i
}

type meters struct {
	duration, active, dropped        *instrument
	memUsed, memAlloc, memGoal, goro *instrument
}

func newMeters() meters {
	return meters{
		duration: newInstrument("http.server.request.duration", "s", histogram),
		active:   newInstrument("http.server.active_requests", "{request}", updown),
		dropped:  newInstrument("hgmx.telemetry.dropped", "{item}", counter),
		memUsed:  newInstrument("go.memory.used", "By", updown),
		memAlloc: newInstrument("go.memory.allocated", "By", counter),
		memGoal:  newInstrument("go.memory.gc.goal", "By", updown),
		goro:     newInstrument("go.goroutine.count", "{goroutine}", updown),
	}
}

func (m *meters) all() []*instrument {
	return []*instrument{m.duration, m.active, m.dropped, m.memUsed, m.memAlloc, m.memGoal, m.goro}
}

var overflow = []slog.Attr{slog.Bool("otel.metric.overflow", true)}

func key(dst []byte, attrs []slog.Attr) []byte {
	for _, a := range attrs {
		dst = binary.AppendUvarint(dst, uint64(len(a.Key)))
		dst = append(dst, a.Key...)
		dst = append(dst, byte(a.Value.Kind()))
		var tmp [24]byte
		var v []byte
		switch a.Value.Kind() {
		case slog.KindString:
			dst = binary.AppendUvarint(dst, uint64(len(a.Value.String())))
			dst = append(dst, a.Value.String()...)
			continue
		case slog.KindInt64:
			v = strconv.AppendInt(tmp[:0], a.Value.Int64(), 10)
		case slog.KindBool:
			v = strconv.AppendBool(tmp[:0], a.Value.Bool())
		default:
			v = []byte(a.Value.String())
		}
		dst = binary.AppendUvarint(dst, uint64(len(v)))
		dst = append(dst, v...)
	}
	return dst
}

func (i *instrument) get(attrs []slog.Attr) *series {
	var buf [128]byte
	k := key(buf[:0], attrs)
	if s := i.series[string(k)]; s != nil {
		return s
	}
	if len(i.series) >= maxSeries-1 {
		attrs, k = overflow, key(buf[:0], overflow)
		if s := i.series[string(k)]; s != nil {
			return s
		}
	}
	s := &series{attrs: slices.Clone(attrs)}
	if i.kind == histogram {
		s.buckets = make([]uint64, len(i.bounds)+1)
	}
	i.series[string(k)] = s
	return s
}

func (i *instrument) add(n int64, attrs ...slog.Attr) {
	i.mu.Lock()
	i.get(attrs).value += n
	i.mu.Unlock()
}

func (i *instrument) set(n int64) {
	i.mu.Lock()
	i.get(nil).value = n
	i.mu.Unlock()
}

func (i *instrument) record(v float64, attrs ...slog.Attr) {
	b, _ := slices.BinarySearch(i.bounds, v)
	i.mu.Lock()
	s := i.get(attrs)
	s.count++
	s.sum += v
	s.buckets[b]++
	i.mu.Unlock()
}

type point struct {
	attrs   []slog.Attr
	value   int64
	count   uint64
	sum     float64
	buckets []uint64
}

type metric struct {
	name, unit string
	kind       kind
	bounds     []float64
	points     []point
}

func (i *instrument) collect() (metric, bool) {
	m := metric{name: i.name, unit: i.unit, kind: i.kind, bounds: i.bounds}
	i.mu.Lock()
	defer i.mu.Unlock()
	for _, s := range i.series {
		m.points = append(m.points, point{attrs: s.attrs, value: s.value, count: s.count, sum: s.sum, buckets: slices.Clone(s.buckets)})
	}
	return m, len(m.points) > 0
}

func (m *meters) gather() []metric {
	samples := []metrics.Sample{
		{Name: "/memory/classes/total:bytes"},
		{Name: "/memory/classes/heap/released:bytes"},
		{Name: "/gc/heap/allocs:bytes"},
		{Name: "/gc/heap/goal:bytes"},
		{Name: "/sched/goroutines:goroutines"},
	}
	metrics.Read(samples)
	u := func(i int) int64 {
		if samples[i].Value.Kind() != metrics.KindUint64 {
			return 0
		}
		return int64(samples[i].Value.Uint64())
	}
	m.memUsed.set(u(0) - u(1))
	m.memAlloc.set(u(2))
	m.memGoal.set(u(3))
	m.goro.set(u(4))

	var out []metric
	for _, i := range m.all() {
		if c, ok := i.collect(); ok {
			out = append(out, c)
		}
	}
	return out
}
