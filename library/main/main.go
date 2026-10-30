package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cogikyo/hgmx/library"
	"github.com/cogikyo/hgmx/server"
	"github.com/cogikyo/hgmx/telemetry"
)

var version = "dev"

const csp = "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; connect-src 'self'; base-uri 'none'; frame-ancestors 'none'; form-action 'self'"

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	dir := flag.String("static", "", "serve /static/ from this directory instead of the embedded files")
	flag.Parse()

	p, err := telemetry.FromEnv(slog.NewTextHandler(os.Stderr, nil))
	log := p.Logger()
	slog.SetDefault(log)
	if err != nil {
		// Invalid exporter configuration leaves the site and local logs usable.
		log.Error("telemetry disabled", "err", err)
	}
	if err := run(p, *addr, *dir); err != nil {
		log.Error("exit", "err", err)
		os.Exit(1)
	}
}

func run(p *telemetry.Pipeline, addr, dir string) error {
	site, err := views.New(views.Config{Version: version, Static: dir, Head: p.Beacon().Head()})
	if err != nil {
		return err
	}
	// Site owns health, the beacon, and the application middleware chain.
	h := server.Site{
		Version:   version,
		CSP:       csp,
		Fallback:  http.HandlerFunc(views.InternalError),
		Telemetry: p,
	}.Handler(site)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop)
	log := p.Logger()
	log.Info("listening", "addr", addr, "version", version)
	return server.Run(ctx, server.New(addr, h, log), p.Shutdown)
}
