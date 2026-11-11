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

func checkCmd(skip []string) (code int) {
	log := newLogger(logLevel, os.Stderr)
	s, err := loadSite(".")
	if err != nil {
		log.Error("Failed to read go.mod", slog.String("error", err.Error()))
		return 1
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	k, err := loadKit(s, skip, log)
	if err != nil {
		log.Error("Failed to plan checks", slog.String("error", err.Error()))
		return 1
	}

	type stage struct {
		name string
		run  func() error
	}
	var stages []stage
	if k != nil {
		stages = append(stages, stage{"contracts", func() error { return k.contracts(ctx) }})
	}
	stages = append(stages, stage{"gen", func() error {
		stale, err := staleOutputs(ctx, s, k)
		if err != nil {
			return err
		}
		for _, line := range stale {
			fmt.Fprintf(os.Stdout, "%s; run go tool hgmx gen\n", line)
		}
		if len(stale) > 0 {
			return fmt.Errorf("stale files: %d", len(stale))
		}
		return nil
	}})
	if k != nil {
		stages = append(stages, stage{"pins", func() error {
			bad := k.pins()
			for _, line := range bad {
				fmt.Fprintln(os.Stdout, line)
			}
			if len(bad) > 0 {
				return fmt.Errorf("runtime pins: %d", len(bad))
			}
			return nil
		}})
	}
	stages = append(stages, []stage{
		{"rules", func() error {
			findings, err := runRules(s)
			for _, f := range findings {
				f.print(os.Stdout)
			}
			if err != nil {
				return err
			}
			if len(findings) > 0 {
				return fmt.Errorf("findings: %d", len(findings))
			}
			return nil
		}},
		{"credentials", func() error { return s.checkCredentials(ctx) }},
		{"fmt", func() error {
			changed, err := formatSite(s, false)
			if err != nil {
				return err
			}
			for _, rel := range changed {
				fmt.Fprintf(os.Stdout, "%s: not formatted; run go tool hgmx fmt\n", rel)
			}
			if len(changed) > 0 {
				return fmt.Errorf("unformatted files: %d", len(changed))
			}
			return nil
		}},
		{"vet", s.gocmd("vet", "./...").Run},
		{"build", s.gocmd("build", "./...").Run},
		{"test", s.gocmd("test", "./...").Run},
	}...)
	for _, st := range stages {
		start := time.Now()
		if err := st.run(); err != nil {
			log.Error("Check failed", slog.String("stage", st.name), slog.String("error", err.Error()))
			k.incomplete()
			return 1
		}
		log.Info("Check passed", slog.String("stage", st.name), slog.Duration("took", time.Since(start).Round(time.Millisecond)))
	}
	if k.incomplete() {
		return 1
	}
	return 0
}

func migrateCmd(args []string) (code int) {
	log := newLogger(logLevel, os.Stderr)
	signal.Notify(make(chan os.Signal, 1), os.Interrupt)
	cmd := exec.Command("go", append([]string{"run", serverPkg, "migrate"}, args...)...)
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
