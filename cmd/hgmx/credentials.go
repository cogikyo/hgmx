package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"slices"
	"syscall"
	"time"
)

const credentialsFile = "credentials.json"

var errStaleCredentials = errors.New(credentialsFile + " is stale; run go tool hgmx credentials")

func credentialsCmd() int {
	log := newLogger(logLevel, os.Stderr)
	s, err := loadSite(".")
	if err != nil {
		log.Error("Failed to read go.mod", slog.String("error", err.Error()))
		return 1
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	err = s.withBinary(func(bin string) error {
		data, err := s.manifest(ctx, bin)
		if err != nil {
			return err
		}
		if old, _ := os.ReadFile(s.path(credentialsFile)); bytes.Equal(old, data) {
			log.Info("Credentials unchanged", slog.String("file", credentialsFile))
			return nil
		}
		if err := writeAtomic(s.path(credentialsFile), data); err != nil {
			return err
		}
		log.Info("Credentials written", slog.String("file", credentialsFile))
		return nil
	})
	if err != nil {
		log.Error("Failed to generate credentials", slog.String("error", err.Error()))
		return 1
	}
	return 0
}

func credentialsCheckCmd(names, flags []string) int {
	log := newLogger(logLevel, os.Stderr)
	s, err := loadSite(".")
	if err != nil {
		log.Error("Failed to read go.mod", slog.String("error", err.Error()))
		return 1
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	code := 0
	err = s.withBinary(func(bin string) error {
		if err := s.currentCredentials(ctx, bin); err != nil {
			return err
		}
		cmd := exec.CommandContext(ctx, bin, slices.Concat([]string{"credentials", "check"}, names, []string{"--"}, flags)...)
		cmd.Dir = s.dir
		cmd.Stdout, cmd.Stderr = os.Stdout, os.Stderr
		err := cmd.Run()
		var exit *exec.ExitError
		if errors.As(err, &exit) {
			code = exit.ExitCode()
			return nil
		}
		return err
	})
	if err != nil {
		log.Error("Failed to check credentials", slog.String("error", err.Error()))
		return 1
	}
	return code
}

func (s site) checkCredentials(ctx context.Context) error {
	if s.mod.Module.Mod.Path == hgmxPath {
		return nil
	}
	return s.withBinary(func(bin string) error { return s.currentCredentials(ctx, bin) })
}

func (s site) currentCredentials(ctx context.Context, bin string) error {
	want, err := s.manifest(ctx, bin)
	if err != nil {
		return err
	}
	have, err := os.ReadFile(s.path(credentialsFile))
	if errors.Is(err, fs.ErrNotExist) || err == nil && !bytes.Equal(have, want) {
		return errStaleCredentials
	}
	return err
}

func (s site) withBinary(fn func(bin string) error) error {
	dir, err := os.MkdirTemp("", "hgmx-credentials-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(dir)
	bin := filepath.Join(dir, "site")
	build := s.gocmd("build", "-o", bin, serverPkg)
	build.Stdout = os.Stderr
	if err := build.Run(); err != nil {
		return fmt.Errorf("build site: %w", err)
	}
	return fn(bin)
}

func (s site) manifest(ctx context.Context, bin string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bin, "credentials", "manifest")
	cmd.Dir = s.dir
	cmd.Env = []string{}
	for _, key := range []string{"PATH", "HOME", "TMPDIR"} {
		if v, ok := os.LookupEnv(key); ok {
			cmd.Env = append(cmd.Env, key+"="+v)
		}
	}
	cmd.WaitDelay = time.Second
	var out bytes.Buffer
	cmd.Stdout, cmd.Stderr = &out, os.Stderr
	err := cmd.Run()
	switch {
	case errors.Is(ctx.Err(), context.DeadlineExceeded):
		return nil, errors.New("site did not exit from `credentials manifest`; dispatch it to credential.Command in main")
	case err != nil:
		os.Stderr.Write(out.Bytes())
		return nil, fmt.Errorf("site `credentials manifest`: %w", err)
	case !json.Valid(out.Bytes()):
		return nil, errors.New("site `credentials manifest` printed something other than JSON")
	}
	return out.Bytes(), nil
}
