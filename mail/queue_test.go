package mail

import (
	"context"
	"errors"
	"log/slog"
	"slices"
	"testing"
	"time"
)

type sendFunc func(context.Context, Message) error

func (f sendFunc) Send(ctx context.Context, m Message) error { return f(ctx, m) }

func TestQueueDrain(t *testing.T) {
	entered, release := make(chan struct{}), make(chan struct{})
	var sent []string
	q := NewQueue(sendFunc(func(ctx context.Context, m Message) error {
		if len(sent) == 0 {
			close(entered)
			select {
			case <-release:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		sent = append(sent, m.Text)
		return nil
	}), 1, slog.New(slog.DiscardHandler))
	t.Cleanup(func() {
		ctx, cancel := context.WithCancel(context.Background())
		cancel()
		q.Close(ctx)
	})
	q.Enqueue(Message{Text: "first"})
	select {
	case <-entered:
	case <-time.After(time.Second):
		t.Fatal("sender did not start")
	}
	q.Enqueue(Message{Text: "second"})
	enqueued := make(chan struct{})
	go func() {
		q.Enqueue(Message{Text: "dropped"})
		close(enqueued)
	}()
	select {
	case <-enqueued:
	case <-time.After(time.Second):
		t.Fatal("Enqueue blocked on a full queue")
	}
	closed := make(chan error, 1)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	go func() { closed <- q.Close(ctx) }()
	select {
	case err := <-closed:
		t.Fatalf("Close returned before drain: %v", err)
	case <-time.After(20 * time.Millisecond):
	}
	close(release)
	if err := <-closed; err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(sent, []string{"first", "second"}) {
		t.Fatalf("sent = %v", sent)
	}
	q.Enqueue(Message{Text: "closed"})
	if err := q.Close(ctx); err != nil {
		t.Fatalf("second Close = %v", err)
	}
}

func TestQueueCloseDeadline(t *testing.T) {
	entered, stopped := make(chan struct{}), make(chan struct{})
	q := NewQueue(sendFunc(func(ctx context.Context, _ Message) error {
		close(entered)
		<-ctx.Done()
		close(stopped)
		return ctx.Err()
	}), 1, slog.New(slog.DiscardHandler))
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	t.Cleanup(func() { q.Close(ctx) })
	q.Enqueue(Message{})
	select {
	case <-entered:
	case <-ctx.Done():
		t.Fatal("sender did not start")
	}
	q.Enqueue(Message{})
	drain, stop := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer stop()
	if err := q.Close(drain); !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("Close = %v, want deadline exceeded", err)
	}
	select {
	case <-stopped:
	case <-ctx.Done():
		t.Fatal("worker context was not canceled")
	}
	if err := q.Close(ctx); err != nil {
		t.Fatalf("Close after cancellation = %v", err)
	}
}
