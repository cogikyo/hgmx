// Package server assembles and runs HTTP sites behind Cloudflare.
package server

import (
	"cmp"
	"context"
	"errors"
	"io"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/cogikyo/hgmx/telemetry"
)

var drainTimeout = 10 * time.Second

const afterTimeout = 5 * time.Second // Fresh cleanup budget after the HTTP drain.

// Site configures the root handler around an application mux.
type Site struct {
	Version       string
	CSP           string
	Allow         []string
	InternalError http.Handler
	Telemetry     *telemetry.Pipeline
}

// Handler wraps mux with telemetry, security headers, and panic recovery.
// The mux type rejects outer wrappers that can hide r.Pattern from telemetry.
// StripPrefix, nested muxes, and route wrappers can still coarsen http.route.
// Cross-origin protection runs last, before mux.
func (s Site) Handler(mux *http.ServeMux) http.Handler {
	if s.InternalError == nil {
		panic("server: Site.InternalError is nil")
	}
	if s.Telemetry == nil {
		panic("server: Site.Telemetry is nil; use telemetry.Disabled for a site without telemetry")
	}
	p := s.Telemetry
	root := http.NewServeMux()
	root.Handle("GET /healthz", health(s.Version))
	if b := p.Beacon(); b.Head() != "" {
		root.Handle(telemetry.Path, b)
	}
	root.Handle("/", p.Handler(recovery(p.Logger(), s.InternalError, http.NewCrossOriginProtection().Handler(mux))))
	return security(s.CSP, s.Allow, root)
}

// New configures a server for use behind Cloudflare.
// Read and write timeouts stay unset to allow uploads and streams.
func New(addr string, h http.Handler, log *slog.Logger) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       2 * time.Minute,
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelWarn),
	}
}

// Run binds synchronously, then serves until failure, shutdown, or ctx cancellation.
// After a successful bind, it drains for up to 10 seconds before the hooks.
// On drain failure it closes connections; handlers may still outlive Run.
// Closing their resources after Run returns can race those handlers.
// The after hooks share a fresh 5-second context and run even if binding fails.
// Run returns all errors joined.
func Run(ctx context.Context, srv *http.Server, after ...func(context.Context) error) error {
	err := serve(ctx, srv)
	hooks, cancel := context.WithTimeout(context.WithoutCancel(ctx), afterTimeout)
	defer cancel()
	for _, f := range after {
		err = errors.Join(err, f(hooks))
	}
	return err
}

func serve(ctx context.Context, srv *http.Server) error {
	ln, err := net.Listen("tcp", cmp.Or(srv.Addr, ":http"))
	if err != nil {
		return err
	}
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ln) }()

	var serr error
	select {
	case serr = <-served:
	case <-ctx.Done():
	}

	drain, cancel := context.WithTimeout(context.WithoutCancel(ctx), drainTimeout)
	defer cancel()
	if err = srv.Shutdown(drain); err != nil {
		err = errors.Join(err, srv.Close())
	}
	if serr == nil {
		serr = <-served
	}
	if !errors.Is(serr, http.ErrServerClosed) {
		err = errors.Join(serr, err)
	}
	return err
}

func health(version string) http.Handler {
	body := "ok " + version
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		io.WriteString(w, body)
	})
}
