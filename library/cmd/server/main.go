package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"github.com/cogikyo/hgmx/cli"
	"github.com/cogikyo/hgmx/credential"
	"github.com/cogikyo/hgmx/migrate"
	"github.com/cogikyo/hgmx/pg"
	"github.com/cogikyo/hgmx/server"
	"github.com/cogikyo/hgmx/telemetry"

	"github.com/cogikyo/hgmx/library/routes"
	"github.com/cogikyo/hgmx/library/store"
	"github.com/cogikyo/hgmx/library/views"
)

var version = "dev"

var secrets credential.Registry

func main() {
	addr := flag.String("addr", "127.0.0.1:8080", "listen address")
	dir := flag.String("static", "", "serve /static/ from this directory instead of the embedded files")
	if code, ok := cli.Dispatch(cli.Commands{
		"credentials": credential.Run(&secrets, flag.CommandLine, nil),
		"migrate":     migrate.Run(store.Schema),
	}); ok {
		os.Exit(code)
	}
	flag.Parse()

	p := telemetry.Setup(slog.NewTextHandler(os.Stderr, nil))
	if err := run(p, *addr, *dir); err != nil {
		p.Logger().Error("exit", "err", err)
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
		Version:   version,
		Telemetry: p,
	}.Handler(mux)

	ctx, stop := cli.NotifyContext()
	defer stop()
	log.InfoContext(ctx, "listening", "addr", addr, "version", version)
	return server.Run(ctx, server.New(addr, h, log), p.Shutdown, func(ctx context.Context) error {
		return pg.Close(ctx, db)
	})
}
