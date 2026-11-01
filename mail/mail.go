// Package mail queues messages for background delivery.
// Request handlers must use [Queue.Enqueue] so responses do not wait for mail.
package mail

import (
	"cmp"
	"context"
	"log/slog"
)

// Message holds one recipient and plain-text or multipart content.
type Message struct {
	To      string // To is a single recipient address.
	Subject string // Subject must not contain CR or LF when sent by SMTP.
	Text    string // Text is the plain-text body.
	HTML    string // HTML adds an alternative body when nonempty.
}

// Sender performs delivery for a [Queue].
type Sender interface {
	// Send performs delivery synchronously.
	Send(context.Context, Message) error
}

// Log writes message content to a logger instead of delivering it.
type Log struct {
	Log *slog.Logger // Log defaults to slog.Default when nil.
}

// Send logs the recipient, subject and plain-text body at Info level.
func (l Log) Send(ctx context.Context, m Message) error {
	cmp.Or(l.Log, slog.Default()).InfoContext(ctx, "mail", "to", m.To, "subject", m.Subject, "text", m.Text)
	return nil
}
