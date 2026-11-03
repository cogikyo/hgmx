package auth

import (
	"container/list"
	"net/netip"
	"sync"
	"time"
)

const tracked = 10000

type limiter struct {
	mu   sync.Mutex
	max  int
	per  time.Duration
	keys map[string]*list.Element
	lru  *list.List
}

type bucket struct {
	key  string
	hits []time.Time
}

func newLimiter(max int, per time.Duration) *limiter {
	return &limiter{max: max, per: per, keys: map[string]*list.Element{}, lru: list.New()}
}

func (l *limiter) allow(ip string) bool {
	key := client(ip)
	now := time.Now()
	l.mu.Lock()
	defer l.mu.Unlock()
	for e := l.lru.Front(); e != nil; e = l.lru.Front() {
		b := e.Value.(*bucket)
		if now.Sub(b.hits[len(b.hits)-1]) < l.per {
			break
		}
		l.drop(e)
	}
	e, ok := l.keys[key]
	if !ok {
		if l.lru.Len() >= tracked {
			l.drop(l.lru.Front())
		}
		e = l.lru.PushBack(&bucket{key: key})
		l.keys[key] = e
	}
	b := e.Value.(*bucket)
	for len(b.hits) > 0 && now.Sub(b.hits[0]) >= l.per {
		b.hits = b.hits[1:]
	}
	if len(b.hits) >= l.max {
		return false
	}
	b.hits = append(b.hits, now)
	l.lru.MoveToBack(e)
	return true
}

func (l *limiter) drop(e *list.Element) {
	l.lru.Remove(e)
	delete(l.keys, e.Value.(*bucket).key)
}

func client(ip string) string {
	addr, err := netip.ParseAddr(ip)
	if err != nil {
		return ip
	}
	if addr = addr.Unmap(); addr.Is4() {
		return addr.String()
	}
	p, _ := addr.Prefix(64)
	return p.String()
}
