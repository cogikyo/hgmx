// Package cli handles site subcommands, exit codes and signal contexts.
package cli

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"maps"
	"os"
	"os/signal"
	"slices"
	"strings"
	"syscall"
)

// ErrUsage marks command-line misuse; [Exit] maps it to status 2.
var ErrUsage = errors.New("usage error")

type usage struct{ error }

func (u usage) Unwrap() []error { return []error{u.error, ErrUsage} }

// Usage formats an error like [fmt.Errorf] and marks it with [ErrUsage].
func Usage(format string, a ...any) error {
	return usage{fmt.Errorf(format, a...)}
}

// Exit returns 0 for nil, 2 for [ErrUsage], and 1 for any other error.
func Exit(err error) int {
	switch {
	case err == nil:
		return 0
	case errors.Is(err, ErrUsage):
		return 2
	}
	return 1
}

// Command runs a subcommand with the arguments after its name.
type Command func(ctx context.Context, args []string) error

// Commands maps subcommand names to their commands.
type Commands map[string]Command

// Dispatch runs os.Args[1] and returns its [Exit] code with ok true.
// Absent or "-" args return ok false; unknown names print usage and return 2.
// Register server flags first; exit when ok, otherwise call flag.Parse.
// Commands get [NotifyContext]; usage errors print to stderr and others log.
func Dispatch(cmds Commands) (code int, ok bool) {
	if len(os.Args) < 2 || strings.HasPrefix(os.Args[1], "-") {
		return 0, false
	}
	name := os.Args[1]
	cmd, found := cmds[name]
	if !found {
		err := Usage("unknown command %q; want %s or flags", name, strings.Join(slices.Sorted(maps.Keys(cmds)), ", "))
		fmt.Fprintln(os.Stderr, err)
		return Exit(err), true
	}
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, nil)))
	ctx, stop := NotifyContext()
	defer stop()
	err := cmd(ctx, os.Args[2:])
	switch {
	case errors.Is(err, ErrUsage):
		fmt.Fprintln(os.Stderr, err)
	case err != nil:
		slog.Error(name, "err", err)
	}
	return Exit(err), true
}

// NotifyContext returns a context canceled by the first SIGINT or SIGTERM.
// A second signal gets the default behavior and ends the process.
func NotifyContext() (context.Context, context.CancelFunc) {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	context.AfterFunc(ctx, stop)
	return ctx, stop
}
