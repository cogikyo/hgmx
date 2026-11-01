package mail

import (
	"cmp"
	"context"
	"log/slog"
	"sync"
	"sync/atomic"
)

// Queue sends messages in order on one background worker.
type Queue struct {
	log     *slog.Logger
	pending chan Message
	done    chan struct{}
	cancel  context.CancelFunc
	mu      sync.Mutex
	closed  bool
	dropped atomic.Uint64
}

// NewQueue starts a worker; capacity must be positive or it panics.
// A nil log uses slog.Default; delivery errors are logged without retry.
func NewQueue(sender Sender, capacity int, log *slog.Logger) *Queue {
	if capacity < 1 {
		panic("mail: queue capacity must be positive")
	}
	ctx, cancel := context.WithCancel(context.Background())
	q := &Queue{
		log:     cmp.Or(log, slog.Default()),
		pending: make(chan Message, capacity),
		done:    make(chan struct{}),
		cancel:  cancel,
	}
	go func() {
		defer close(q.done)
		defer cancel()
		flush := func() {
			if n := q.dropped.Swap(0); n > 0 {
				q.log.WarnContext(ctx, "mail queue full; messages dropped", "count", n)
			}
		}
		defer flush()
		for m := range q.pending {
			flush()
			if ctx.Err() != nil {
				q.log.WarnContext(ctx, "mail queue stopped; messages dropped", "count", 1+len(q.pending))
				return
			}
			if err := sender.Send(ctx, m); err != nil {
				q.log.ErrorContext(ctx, "mail send failed", "err", err)
			}
		}
	}()
	return q
}

// Enqueue never waits for delivery or queue space; full or closed queues drop m.
func (q *Queue) Enqueue(m Message) {
	q.mu.Lock()
	if q.closed {
		q.mu.Unlock()
		return
	}
	select {
	case q.pending <- m:
		q.mu.Unlock()
	default:
		q.dropped.Add(1)
		q.mu.Unlock()
	}
}

// Close stops accepting messages and waits for pending delivery.
// If ctx ends first, it cancels delivery and returns ctx.Err() without waiting.
// Pass it as an after hook to [github.com/cogikyo/hgmx/server.Run].
func (q *Queue) Close(ctx context.Context) error {
	q.mu.Lock()
	if !q.closed {
		q.closed = true
		close(q.pending)
	}
	q.mu.Unlock()
	select {
	case <-q.done:
		return nil
	default:
	}
	select {
	case <-q.done:
		return nil
	case <-ctx.Done():
		q.cancel()
		return ctx.Err()
	}
}
