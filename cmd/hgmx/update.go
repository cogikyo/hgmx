package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

const (
	hgmxPath       = "github.com/cogikyo/hgmx"
	hgmxTool       = hgmxPath + "/cmd/hgmx"
	vendorRecovery = "vendor/ may be incomplete or out of sync; recover with GOFLAGS=-mod=mod go mod vendor"
)

func moduleEnv() []string {
	return append(os.Environ(), "GOWORK=off", "GOFLAGS=-mod=mod", "GIT_TERMINAL_PROMPT=0")
}

func latestHgmx(ctx context.Context) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "git", "ls-remote", "https://"+hgmxPath, "refs/heads/master")
	cmd.Env = moduleEnv()
	cmd.WaitDelay = 100 * time.Millisecond
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		if ctx.Err() != nil {
			return "", fmt.Errorf("git ls-remote: %w", ctx.Err())
		}
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			err = fmt.Errorf("%w: %s", err, msg)
		}
		return "", fmt.Errorf("git ls-remote: %w", err)
	}
	fields := strings.Fields(string(out))
	if len(fields) != 2 || len(fields[0]) != 40 || fields[1] != "refs/heads/master" {
		return "", fmt.Errorf("git ls-remote: invalid master response %q", strings.TrimSpace(string(out)))
	}
	if _, err := hex.DecodeString(fields[0]); err != nil {
		return "", fmt.Errorf("git ls-remote: invalid master commit %q: %w", fields[0], err)
	}
	return fields[0][:12], nil
}

func hgmxVersion(ctx context.Context, commit string) (string, error) {
	dir, err := os.MkdirTemp("", "hgmx-version-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(dir)
	cmd := exec.CommandContext(ctx, "go", "list", "-m", "-f", "{{.Version}}", hgmxPath+"@"+commit)
	cmd.Dir, cmd.Env = dir, moduleEnv()
	cmd.WaitDelay = 100 * time.Millisecond
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		if msg := strings.TrimSpace(stderr.String()); msg != "" {
			err = fmt.Errorf("%w: %s", err, msg)
		}
		return "", fmt.Errorf("go list: %w", err)
	}
	v := strings.TrimSpace(string(out))
	if !semver.IsValid(v) {
		return "", fmt.Errorf("invalid hgmx version %q for commit %s", v, commit)
	}
	return v, nil
}

func (s site) version(path string) string {
	for _, r := range s.mod.Require {
		if r.Mod.Path == path {
			return r.Mod.Version
		}
	}
	return ""
}

func (s site) hgmxReplaces() []*modfile.Replace {
	var out []*modfile.Replace
	for _, r := range s.mod.Replace {
		if r.Old.Path == hgmxPath || r.Old.Path == hgmxTool {
			out = append(out, r)
		}
	}
	return out
}

func updateCmd() int {
	log := newLogger(logLevel, os.Stderr)
	s, err := loadSite(".")
	if err != nil {
		log.Error("Failed to read go.mod; run update from the site's module root", slog.String("error", err.Error()))
		return 1
	}
	if s.mod.Module.Mod.Path == hgmxPath {
		log.Error("Refusing to update inside the hgmx module; run update from a site's module root")
		return 1
	}
	replaces := s.hgmxReplaces()
	for _, r := range replaces {
		log.Error("Refusing to update with an hgmx replace; remove this line and use go.work to develop against ../hgmx",
			slog.String("file", "go.mod"), slog.Int("line", r.Syntax.Start.Line), slog.String("replace", strings.Join(r.Syntax.Token, " ")))
	}
	if len(replaces) > 0 {
		return 1
	}
	old, oldTempl := s.version(hgmxPath), s.version(templPath)
	if old == "" {
		old = "(none)"
	}
	originals := []struct {
		name   string
		data   []byte
		exists bool
	}{{name: "go.mod"}, {name: "go.sum"}}
	for i := range originals {
		f := &originals[i]
		f.data, err = os.ReadFile(s.path(f.name))
		f.exists = err == nil
		if err != nil && !(f.name == "go.sum" && errors.Is(err, os.ErrNotExist)) {
			log.Error("Failed to snapshot module file", slog.String("file", f.name), slog.String("error", err.Error()))
			return 1
		}
	}
	restore := func() {
		restored := true
		for _, f := range originals {
			var err error
			if f.exists {
				err = writeAtomic(s.path(f.name), f.data)
			} else {
				err = os.Remove(s.path(f.name))
				if errors.Is(err, os.ErrNotExist) {
					err = nil
				}
			}
			if err != nil {
				restored = false
				log.Error("Failed to restore module file; restore it before retrying update", slog.String("file", f.name), slog.String("error", err.Error()))
			}
		}
		if restored {
			log.Info("Restored go.mod and go.sum to their pre-update state")
		}
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	latest, err := latestHgmx(ctx)
	if err != nil {
		log.Error("Could not resolve hgmx master; update aborted", slog.String("error", err.Error()))
		return 1
	}
	target, err := hgmxVersion(ctx, latest)
	if err != nil {
		log.Error("Could not resolve hgmx master commit's module version; update aborted", slog.String("commit", latest), slog.String("error", err.Error()))
		return 1
	}
	if semver.Compare(target, s.version(hgmxPath)) < 0 {
		log.Error("Refusing to downgrade hgmx", slog.String("current", old), slog.String("target", target), slog.String("master", latest))
		return 1
	}
	for _, args := range [][]string{
		{"get", "-tool", hgmxTool + "@" + latest},
		{"mod", "tidy"},
		{"mod", "vendor"},
	} {
		log.Info("Running go " + strings.Join(args, " "))
		cmd := exec.CommandContext(ctx, "go", args...)
		cmd.Dir, cmd.Env = s.dir, moduleEnv()
		var output bytes.Buffer
		cmd.Stdin = os.Stdin
		cmd.Stdout = io.MultiWriter(os.Stderr, &output)
		cmd.Stderr = cmd.Stdout
		if err := cmd.Run(); err != nil {
			restore()
			log.Error("Update failed", slog.String("command", "go "+strings.Join(args, " ")), slog.String("error", err.Error()), slog.String("output", strings.TrimSpace(output.String())))
			if args[0] == "mod" && args[1] == "vendor" {
				log.Warn(vendorRecovery)
			}
			return 1
		}
	}
	updated, err := loadSite(s.dir)
	if err != nil {
		restore()
		log.Error("Failed to read updated go.mod", slog.String("error", err.Error()))
		log.Warn(vendorRecovery)
		return 1
	}
	s = updated
	log.Info("hgmx " + old + " → " + s.version(hgmxPath))
	next := "go tool hgmx check"
	if oldTempl != s.version(templPath) {
		next = "go tool hgmx gen, then " + next
	}
	log.Info("Next: " + next)
	return 0
}
