package mail

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"mime/quotedprintable"
	"net"
	stdmail "net/mail"
	"net/smtp"
	"net/textproto"
	"strconv"
	"strings"
	"time"

	"github.com/cogikyo/hgmx/credential"
)

// SMTP delivers mail over implicit TLS, not STARTTLS.
type SMTP struct {
	Host     string // Host is the server name used for TLS verification.
	Port     int    // Port defaults to 465 when zero.
	User     string // User is the SMTP authentication name.
	Password string // Password is the SMTP authentication secret.
	From     string // From is the sender address, with an optional display name.
}

// Send delivers synchronously with a timeout of at most 30 seconds.
func (s SMTP) Send(ctx context.Context, m Message) error {
	from, err := stdmail.ParseAddress(s.From)
	if err != nil {
		return fmt.Errorf("mail: invalid from address: %w", err)
	}
	to, err := stdmail.ParseAddress(m.To)
	if err != nil {
		return fmt.Errorf("mail: invalid to address: %w", err)
	}
	body, err := message(from, to, m)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	client, err := s.dial(ctx)
	if err != nil {
		return err
	}
	defer client.Close()
	if err := client.Auth(smtp.PlainAuth("", s.User, s.Password, s.Host)); err != nil {
		return fmt.Errorf("mail: SMTP auth: %w", err)
	}
	if err := client.Mail(from.Address); err != nil {
		return fmt.Errorf("mail: SMTP from: %w", err)
	}
	if err := client.Rcpt(to.Address); err != nil {
		return fmt.Errorf("mail: SMTP to: %w", err)
	}
	w, err := client.Data()
	if err != nil {
		return fmt.Errorf("mail: SMTP data: %w", err)
	}
	if _, err := w.Write(body); err != nil {
		return fmt.Errorf("mail: SMTP body: %w", err)
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("mail: SMTP delivery: %w", err)
	}
	if err := client.Quit(); err != nil {
		return fmt.Errorf("mail: SMTP quit: %w", err)
	}
	return nil
}

// Check authenticates without sending mail; it does not verify delivery.
func (s SMTP) Check(ctx context.Context) error {
	if s.Host == "" || s.User == "" || s.Password == "" {
		return &credential.Failure{Reason: "SMTP host, user and password must all be set"}
	}
	if _, err := stdmail.ParseAddress(s.From); err != nil {
		return &credential.Failure{Reason: "SMTP from address is invalid"}
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	client, err := s.dial(ctx)
	var cert *tls.CertificateVerificationError
	var timeout net.Error
	switch {
	case errors.As(err, &cert):
		return &credential.Failure{Reason: "SMTP server certificate not trusted for " + s.Host}
	case errors.As(err, &timeout) && timeout.Timeout():
		return &credential.Failure{Reason: "SMTP connection timed out"}
	case err != nil:
		return &credential.Failure{Reason: "SMTP connection failed"}
	}
	defer client.Close()
	if ok, _ := client.Extension("AUTH"); !ok {
		return &credential.Failure{Reason: "SMTP server offers no AUTH"}
	}
	if err := client.Auth(smtp.PlainAuth("", s.User, s.Password, s.Host)); err != nil {
		var reply *textproto.Error
		if errors.As(err, &reply) {
			return &credential.Failure{Reason: "SMTP authentication rejected", Code: reply.Code}
		}
		return &credential.Failure{Reason: "SMTP authentication failed"}
	}
	if err := client.Quit(); err != nil {
		return &credential.Failure{Reason: "SMTP quit failed after authentication"}
	}
	return nil
}

func (s SMTP) dial(ctx context.Context) (*smtp.Client, error) {
	port := s.Port
	if port == 0 {
		port = 465
	}
	dialer := tls.Dialer{
		NetDialer: &net.Dialer{Timeout: 10 * time.Second},
		Config:    &tls.Config{ServerName: s.Host, MinVersion: tls.VersionTLS12},
	}
	conn, err := dialer.DialContext(ctx, "tcp", net.JoinHostPort(s.Host, strconv.Itoa(port)))
	if err != nil {
		return nil, fmt.Errorf("mail: SMTP connect: %w", err)
	}
	context.AfterFunc(ctx, func() { conn.Close() })
	deadline, _ := ctx.Deadline()
	if err := conn.SetDeadline(deadline); err != nil {
		conn.Close()
		return nil, fmt.Errorf("mail: SMTP deadline: %w", err)
	}
	client, err := smtp.NewClient(conn, s.Host)
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("mail: SMTP greeting: %w", err)
	}
	return client, nil
}

func message(from, to *stdmail.Address, m Message) ([]byte, error) {
	if strings.ContainsAny(m.Subject, "\r\n") {
		return nil, fmt.Errorf("mail: subject contains a newline")
	}
	domain := from.Address[strings.LastIndexByte(from.Address, '@')+1:]
	var b bytes.Buffer
	fmt.Fprintf(&b, "From: %s\r\nTo: %s\r\nSubject: %s\r\nDate: %s\r\nMessage-ID: <%s@%s>\r\nMIME-Version: 1.0\r\n",
		from.String(), to.String(), mime.QEncoding.Encode("utf-8", m.Subject), time.Now().Format(time.RFC1123Z), rand.Text(), domain)
	if m.HTML == "" {
		b.WriteString("Content-Type: text/plain; charset=utf-8\r\nContent-Transfer-Encoding: quoted-printable\r\n\r\n")
		w := quotedprintable.NewWriter(&b)
		if _, err := io.WriteString(w, m.Text); err != nil {
			return nil, fmt.Errorf("mail: encode text: %w", err)
		}
		if err := w.Close(); err != nil {
			return nil, fmt.Errorf("mail: encode text: %w", err)
		}
		return b.Bytes(), nil
	}
	w := multipart.NewWriter(&b)
	fmt.Fprintf(&b, "Content-Type: multipart/alternative; boundary=%q\r\n\r\n", w.Boundary())
	for _, part := range []struct{ kind, body string }{{"text/plain", m.Text}, {"text/html", m.HTML}} {
		p, err := w.CreatePart(textproto.MIMEHeader{
			"Content-Type":              {part.kind + "; charset=utf-8"},
			"Content-Transfer-Encoding": {"quoted-printable"},
		})
		if err != nil {
			return nil, fmt.Errorf("mail: MIME part: %w", err)
		}
		qp := quotedprintable.NewWriter(p)
		if _, err := io.WriteString(qp, part.body); err != nil {
			return nil, fmt.Errorf("mail: encode %s: %w", part.kind, err)
		}
		if err := qp.Close(); err != nil {
			return nil, fmt.Errorf("mail: encode %s: %w", part.kind, err)
		}
	}
	if err := w.Close(); err != nil {
		return nil, fmt.Errorf("mail: MIME close: %w", err)
	}
	return b.Bytes(), nil
}
