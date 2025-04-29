package sloghandler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"strings"
	"sync"

	"github.com/fatih/color"
)

var _ slog.Handler = &Handler{}

type Handler struct {
	h slog.Handler
	m *sync.Mutex
	w io.Writer
}

var levelToIcon = map[slog.Level]string{
	slog.LevelDebug: "(?)",
	slog.LevelInfo:  "(i)",
	slog.LevelWarn:  "(!)",
	slog.LevelError: "(x)",
}

var levelToColor = map[slog.Level]*color.Color{
	slog.LevelDebug: color.New(color.FgMagenta),
	slog.LevelInfo:  color.New(color.FgBlue),
	slog.LevelWarn:  color.New(color.FgYellow),
	slog.LevelError: color.New(color.FgRed),
}

var attributeColor = color.New(color.Faint)

func NewHandler(w io.Writer, opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	minLevel := slog.LevelDebug
	if opts.Level != nil {
		minLevel = opts.Level.Level()
	}

	return &Handler{
		w: w,
		h: slog.NewTextHandler(io.Discard, &slog.HandlerOptions{
			Level:     minLevel,
			AddSource: opts.AddSource,
		}),
		m: &sync.Mutex{},
	}
}

func NewLogger(logLevel string, verbose bool, stderr io.Writer) *slog.Logger {
	level := slog.LevelInfo
	switch logLevel {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		if logLevel != "info" {
			fmt.Fprintf(stderr, "Invalid log level %q, defaulting to info\n", logLevel)
		}
	}

	if verbose {
		logLevel = "debug"
	}

	handler := NewHandler(stderr, &slog.HandlerOptions{Level: level})
	return slog.New(handler)
}

func Pretty(v any) string {
	b, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return fmt.Sprintf("%v", v)
	}
	return string(b)
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.h.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{h: h.h.WithAttrs(attrs), w: h.w, m: h.m}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{h: h.h.WithGroup(name), w: h.w, m: h.m}
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) (err error) {
	if !h.Enabled(ctx, r.Level) {
		return nil
	}

	icon, okIcon := levelToIcon[r.Level]
	clr, okColor := levelToColor[r.Level]

	if !okIcon {
		icon = "(_)"
	}
	if !okColor {
		clr = color.New()
	}

	var sb strings.Builder
	sb.WriteString(clr.Sprintf("%s %s", icon, r.Message))

	numAttrs := r.NumAttrs()
	switch {
	case numAttrs == 0:
		sb.WriteString("\n")
	case numAttrs <= 2:
		sb.WriteString(" ")
		first := true
		r.Attrs(func(a slog.Attr) bool {
			if !first {
				sb.WriteString(" ")
			}
			first = false
			sb.WriteString(attributeColor.Sprintf("%s=\"%v\"", a.Key, a.Value.Any()))
			return true
		})
		sb.WriteString("\n")
	case numAttrs > 2:
		sb.WriteString("\n")
		r.Attrs(func(a slog.Attr) bool {
			sb.WriteString("    ")
			sb.WriteString(attributeColor.Sprintf("%s=\"%v\"", a.Key, a.Value.Any()))
			sb.WriteString("\n")
			return true
		})
	}
	h.m.Lock()
	defer h.m.Unlock()

	output := sb.String()
	_, err = io.WriteString(h.w, output)
	return err
}
