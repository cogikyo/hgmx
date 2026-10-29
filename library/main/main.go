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

	local := slog.NewTextHandler(os.Stderr, nil)
	if err := run(local, *addr, *dir); err != nil {
		slog.New(local).Error("exit", "err", err)
		os.Exit(1)
	}
}

func run(local slog.Handler, addr, dir string) error {
	p, err := telemetry.FromEnv()
	log := slog.New(p.Logs(local))
	slog.SetDefault(log)
	if err != nil {
		log.Error("telemetry disabled", "err", err)
	}
	beacon, err := p.Beacon()
	if err != nil {
		return err
	}

	site := http.NewServeMux()
	if err := views.Mount(site, version, dir, beacon.Head()); err != nil {
		return err
	}

	root := http.NewServeMux()
	root.Handle("GET /healthz", server.Health(version))
	root.Handle(telemetry.Path, http.NewCrossOriginProtection().Handler(beacon))
	root.Handle("/", p.Handler(local, server.Security(csp, server.Recover(log, http.HandlerFunc(views.InternalError), http.NewCrossOriginProtection().Handler(site)))))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop)
	log.Info("listening", "addr", addr, "version", version)
	return server.Run(ctx, server.New(addr, root, log), p.Shutdown)
}
