package main

import (
	"context"
	"encoding/base64"
	"flag"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"

	"github.com/a-h/templ"

	"github.com/cogikyo/hgmx/auth"
	"github.com/cogikyo/hgmx/credential"
	"github.com/cogikyo/hgmx/mail"
	"github.com/cogikyo/hgmx/migrate"
	"github.com/cogikyo/hgmx/pg"
	"github.com/cogikyo/hgmx/server"
	"github.com/cogikyo/hgmx/telemetry"

	"example.com/fixture/api"
)

var (
	secrets credential.Registry
	authKey = secrets.Declare(auth.KeyCredential, credential.Always, "site")
)

func main() {
	addr := flag.String("addr", "127.0.0.1:18550", "listen address")
	origin := flag.String("origin", "http://localhost:18550", "public origin")
	signer := flag.String("android-signer", "", "base64 SHA-256 of the Android signing certificate")
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "credentials":
			os.Exit(credential.Command(context.Background(), &secrets, flag.CommandLine, os.Args[2:], nil, os.Stdout))
		case "migrate":
			os.Exit(migrateMain(os.Args[2:]))
		}
	}
	flag.Parse()
	p, err := telemetry.FromEnv(slog.NewTextHandler(os.Stderr, nil))
	log := p.Logger()
	if err != nil {
		log.Error("telemetry disabled", "err", err)
	}
	if err := run(p, *addr, *origin, *signer); err != nil {
		log.Error("exit", "err", err)
		os.Exit(1)
	}
}

func run(p *telemetry.Pipeline, addr, origin, signer string) error {
	log := p.Logger()
	key, err := authKey.Read()
	if err != nil {
		return err
	}
	cert, err := base64.StdEncoding.DecodeString(signer)
	if err != nil {
		return err
	}
	db, err := pg.Open(log)
	if err != nil {
		return err
	}
	q := mail.NewQueue(mail.Log{Log: log}, 16, log)
	layout := func(_ *http.Request, _ string, body templ.Component) templ.Component { return body }
	hosts := []auth.Host{{
		Origin: origin,
		Mode:   auth.Open,
		Layout: layout,
		Apps: []auth.App{{
			ID:      "fixture-android",
			Android: &auth.Android{Package: "com.example.fixture.notes", Signers: [][]byte{cert}},
			Build:   1,
			Store:   "https://play.google.com/store/apps/details?id=com.example.fixture.notes",
		}, {
			ID:    "fixture-ios",
			IOS:   &auth.IOS{Prefix: "FIXTURE000", Bundle: "com.example.fixture.notes"},
			Build: 1,
			Store: "https://apps.apple.com/app/id0000000000",
		}},
	}}
	if u, err := url.Parse(origin); err == nil && u.Port() != "" {
		hosts = append(hosts, auth.Host{Origin: "https://" + u.Hostname(), Mode: auth.Open, Layout: layout})
	}
	a, err := auth.New(db, q, auth.Config{Name: "Fixture", Key: []byte(key), Hosts: hosts})
	if err != nil {
		return err
	}
	mux := http.NewServeMux()
	a.Mount(mux)
	rpc := a.API()
	api.Register(rpc)
	if err := rpc.Mount(mux); err != nil {
		return err
	}
	h := server.Site{Telemetry: p, RPC: rpc.Server()}.Handler(mux)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	return server.Run(ctx, server.New(addr, h, log), q.Close, p.Shutdown)
}

func migrateMain(args []string) int {
	log := slog.New(slog.NewTextHandler(os.Stderr, nil))
	db, err := pg.Open(log)
	if err == nil {
		err = migrate.Command(context.Background(), db, []migrate.Track{auth.Track}, args, os.Stdout)
	}
	if err != nil {
		log.Error("migrate", "err", err)
		return 1
	}
	return 0
}
