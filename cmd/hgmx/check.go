package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"
)

func checkCmd() (code int) {
	log := newLogger(logLevel, os.Stderr)
	s, err := loadSite(".")
	if err != nil {
		log.Error("Failed to read go.mod", slog.String("error", err.Error()))
		return 1
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	stages := []struct {
		name string
		run  func() error
	}{
		{"rules", func() error {
			findings, err := runRules(s)
			if err != nil {
				return err
			}
			for _, f := range findings {
				f.print(os.Stdout)
			}
			if len(findings) > 0 {
				return fmt.Errorf("findings: %d", len(findings))
			}
			return nil
		}},
		{"gen", func() error {
			stale, err := genCheck(ctx, s)
			if err != nil {
				return err
			}
			printStale(os.Stdout, stale, "hgmx gen")
			if len(stale) > 0 {
				return fmt.Errorf("stale files: %d", len(stale))
			}
			return nil
		}},
		{"fmt", func() error {
			changed, err := formatSite(s, false)
			if err != nil {
				return err
			}
			for _, rel := range changed {
				fmt.Fprintf(os.Stdout, "%s: not formatted; run hgmx fmt\n", rel)
			}
			if len(changed) > 0 {
				return fmt.Errorf("unformatted files: %d", len(changed))
			}
			return nil
		}},
		{"vet", s.gocmd("vet", "./...").Run},
		{"build", s.gocmd("build", "./...").Run},
		{"test", s.gocmd("test", "./...").Run},
	}
	for _, st := range stages {
		start := time.Now()
		if err := st.run(); err != nil {
			log.Error("Check failed", slog.String("stage", st.name), slog.String("error", err.Error()))
			return 1
		}
		log.Info("Check passed", slog.String("stage", st.name), slog.Duration("took", time.Since(start).Round(time.Millisecond)))
	}
	return 0
}

func migrateCmd(args []string) (code int) {
	log := newLogger(logLevel, os.Stderr)
	signal.Notify(make(chan os.Signal, 1), os.Interrupt)
	cmd := exec.Command("go", append([]string{"run", ".", "migrate"}, args...)...)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	err := cmd.Run()
	var exit *exec.ExitError
	if errors.As(err, &exit) {
		return exit.ExitCode()
	}
	if err != nil {
		log.Error("Failed to run migrate", slog.String("error", err.Error()))
		return 1
	}
	return 0
}
