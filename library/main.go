package main

import (
	"context"
	"flag"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/cogikyo/hgmx/migrate"
	"github.com/cogikyo/hgmx/pg"
	"github.com/cogikyo/hgmx/server"
	"github.com/cogikyo/hgmx/telemetry"

	"github.com/cogikyo/hgmx/library/routes"
	"github.com/cogikyo/hgmx/library/store"
	"github.com/cogikyo/hgmx/library/store/migrations"
	"github.com/cogikyo/hgmx/library/views"
)

var version = "dev"

const csp = "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self'; font-src 'self'; connect-src 'self'; base-uri 'none'; frame-ancestors 'none'; form-action 'self'"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "migrate" {
		os.Exit(migrateMain(os.Args[2:]))
	}

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
	log := p.Logger()
	db, err := pg.Open(log)
	if err != nil {
		return err
	}
	v, err := views.New(views.Config{Version: version, Static: dir, Head: p.Beacon().Head()})
	if err != nil {
		return err
	}
	mux := routes.New(store.New(db), v)
	// Site owns health, the beacon, and the application middleware chain.
	h := server.Site{
		Version:       version,
		CSP:           csp,
		InternalError: http.HandlerFunc(routes.InternalError),
		Telemetry:     p,
	}.Handler(mux)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	context.AfterFunc(ctx, stop)
	log.InfoContext(ctx, "listening", "addr", addr, "version", version)
	return server.Run(ctx, server.New(addr, h, log), p.Shutdown, func(ctx context.Context) error {
		return pg.Close(ctx, db)
	})
}

func migrateMain(args []string) int {
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := pg.Open(log)
	if err == nil {
		err = migrate.Command(ctx, db, []migrate.Track{{Steps: migrations.All}}, args, os.Stdout)
	}
	if err != nil {
		log.ErrorContext(ctx, "migrate", "err", err)
		return 1
	}
	return 0
}
