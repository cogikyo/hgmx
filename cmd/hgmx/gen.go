package main

import (
	"bytes"
	"context"
	"fmt"
	"go/format"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/a-h/templ/cmd/templ/fmtcmd"
	"github.com/a-h/templ/cmd/templ/generatecmd"
)

func genCmd(skip []string) (code int) {
	log := newLogger(logLevel, os.Stderr)
	s, err := loadSite(".")
	if err != nil {
		log.Error("Failed to read go.mod", slog.String("error", err.Error()))
		return 1
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if s.templ() != nil {
		if err := generatecmd.Run(ctx, os.Stdout, os.Stderr, []string{"-path", s.dir}); err != nil {
			log.Error("templ generate failed", slog.String("error", err.Error()))
			return 1
		}
	}
	if s.css() {
		cmd := tailwind(s, s.path(cssOutput))
		cmd.Stdout, cmd.Stderr = os.Stderr, os.Stderr
		if err := cmd.Run(); err != nil {
			log.Error("tailwindcss failed", slog.String("error", err.Error()))
			return 1
		}
	}
	k, err := loadKit(s, skip, log)
	if err != nil {
		log.Error("Failed to plan generation", slog.String("error", err.Error()))
		return 1
	}
	if k != nil {
		if err := k.gen(ctx); err != nil {
			log.Error("Contract generation failed", slog.String("error", err.Error()))
			return 1
		}
	}
	return 0
}

func tailwind(s site, out string, args ...string) *exec.Cmd {
	cmd := exec.Command("tailwindcss", append([]string{"-i", s.path(cssInput), "-o", out, "--minify"}, args...)...)
	cmd.Dir = s.dir
	return cmd
}

func staleOutputs(ctx context.Context, s site, k *kit) ([]string, error) {
	var stale []string
	if s.templ() != nil {
		args, log, _, err := generatecmd.NewArguments(io.Discard, os.Stderr, []string{"-path", s.dir, "-keep-orphaned-files", "-log-level", "warn"})
		if err != nil {
			return nil, err
		}
		writer, changed := generatecmd.NewCheckWriter()
		args.FileWriter = writer
		g, err := generatecmd.NewGenerate(log, args)
		if err != nil {
			return nil, err
		}
		if err := g.Run(ctx); err != nil {
			return nil, err
		}
		for _, name := range changed() {
			rel, err := filepath.Rel(s.dir, name)
			if err != nil {
				return nil, err
			}
			stale = append(stale, filepath.ToSlash(rel)+": out of date")
		}
		err = s.walk(func(rel string, d fs.DirEntry) error {
			src, ok := strings.CutSuffix(rel, "_templ.go")
			if !ok || d.IsDir() {
				return nil
			}
			if _, err := os.Stat(s.path(src + ".templ")); os.IsNotExist(err) {
				stale = append(stale, rel+": no "+src+".templ")
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	}

	if s.css() {
		tmp, err := os.MkdirTemp("", "hgmx-css-")
		if err != nil {
			return nil, err
		}
		defer os.RemoveAll(tmp)
		out := filepath.Join(tmp, "main.min.css")
		cmd := tailwind(s, out)
		if msg, err := cmd.CombinedOutput(); err != nil {
			return nil, fmt.Errorf("tailwindcss: %w\n%s", err, msg)
		}
		want, err := os.ReadFile(out)
		if err != nil {
			return nil, err
		}
		if got, err := os.ReadFile(s.path(cssOutput)); err != nil || !bytes.Equal(got, want) {
			stale = append(stale, cssOutput+": out of date")
		}
	}
	if k != nil {
		lines, err := k.stale(ctx)
		if err != nil {
			return nil, err
		}
		stale = append(stale, lines...)
	}
	return stale, nil
}

func fmtCmd() (code int) {
	log := newLogger(logLevel, os.Stderr)
	s, err := loadSite(".")
	if err != nil {
		log.Error("Failed to read go.mod", slog.String("error", err.Error()))
		return 1
	}
	changed, err := formatSite(s, true)
	if err != nil {
		log.Error("Format failed", slog.String("error", err.Error()))
		return 1
	}
	for _, rel := range changed {
		log.Info("Formatted", slog.String("file", rel))
	}
	return 0
}

func formatSite(s site, write bool) ([]string, error) {
	quiet := slog.New(slog.DiscardHandler)
	var changed []string
	err := s.walk(func(rel string, d fs.DirEntry) error {
		goFile := strings.HasSuffix(rel, ".go")
		if d.IsDir() || !goFile && !strings.HasSuffix(rel, ".templ") {
			return nil
		}
		path := s.path(rel)
		src, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		var formatted []byte
		if goFile {
			formatted, err = format.Source(src)
		} else {
			var buf bytes.Buffer
			err = fmtcmd.Run(quiet, nil, &buf, fmtcmd.Arguments{ToStdout: true, Files: []string{path}, WorkerCount: 1})
			formatted = buf.Bytes()
		}
		if err != nil {
			return fmt.Errorf("%s: %w", rel, err)
		}
		if bytes.Equal(src, formatted) {
			return nil
		}
		changed = append(changed, rel)
		if write {
			return writeAtomic(path, formatted)
		}
		return nil
	})
	return changed, err
}

func devCmd(args []string) (code int) {
	log := newLogger(logLevel, os.Stderr)
	s, err := loadSite(".")
	if err != nil {
		log.Error("Failed to read go.mod", slog.String("error", err.Error()))
		return 1
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var tw *exec.Cmd
	exited := make(chan error, 1)
	if s.css() {
		tw = tailwind(s, s.path(cssOutput), "--watch=always")
		tw.Stdout, tw.Stderr = os.Stderr, os.Stderr
		ownGroup(tw)
		if err := tw.Start(); err != nil {
			log.Error("Failed to start tailwindcss", slog.String("error", err.Error()))
			return 1
		}
		go func() {
			err := tw.Wait()
			if ctx.Err() == nil {
				err = fmt.Errorf("exited: %v", err)
			} else {
				err = nil
			}
			cancel()
			exited <- err
		}()
	}

	run := "go run . -static views/static"
	for _, a := range args {
		run += " '" + strings.ReplaceAll(a, "'", `'\''`) + "'"
	}
	err = generatecmd.Run(ctx, os.Stdout, os.Stderr, []string{"-path", s.dir, "-watch", "-open-browser=false", "-cmd", run})
	failed := err != nil && ctx.Err() == nil
	if failed {
		log.Error("templ watch failed", slog.String("error", err.Error()))
	}
	cancel()
	if tw != nil {
		killGroup(tw.Process)
		if err := <-exited; err != nil {
			log.Error("tailwindcss failed", slog.String("error", err.Error()))
			failed = true
		}
	}
	if failed {
		return 1
	}
	return 0
}
