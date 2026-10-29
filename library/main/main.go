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
)

var version = "dev"

const csp = "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; connect-src 'self'; base-uri 'none'; frame-ancestors 'none'; form-action 'self'"

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	dir := flag.String("static", "", "serve /static/ from this directory instead of the embedded files")
	flag.Parse()

	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	slog.SetDefault(log)
	if err := run(log, *addr, *dir); err != nil {
		log.Error("exit", "err", err)
		os.Exit(1)
	}
}

func run(log *slog.Logger, addr, dir string) error {
	site := http.NewServeMux()
	fallback, err := views.Mount(site, version, dir)
	if err != nil {
		return err
	}

	root := http.NewServeMux()
	root.Handle("GET /healthz", server.Health(version))
	root.Handle("/", server.Log(log, server.Security(csp, server.Recover(log, fallback, http.NewCrossOriginProtection().Handler(site)))))

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	log.Info("listening", "addr", addr, "version", version)
	return server.Run(ctx, server.New(addr, root, log))
}
