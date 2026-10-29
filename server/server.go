// Package server assembles and runs HTTP sites behind Cloudflare.
// Given an application mux, a fallback handler, and a telemetry pipeline:
//
//	h := server.Site{Version: version, CSP: csp, Fallback: fallback, Telemetry: p}.Handler(mux)
//	err := server.Run(ctx, server.New(":8080", h, p.Logger()), p.Shutdown)
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
	Version   string
	CSP       string
	Fallback  http.Handler
	Telemetry *telemetry.Pipeline
}

// Handler wraps mux with telemetry, security headers, and panic recovery.
// Cross-origin protection runs last, before mux.
// It serves /healthz and the telemetry beacon outside that chain.
// It panics if Fallback or Telemetry is nil; use [telemetry.Disabled] to opt out.
func (s Site) Handler(mux http.Handler) http.Handler {
	if s.Fallback == nil {
		panic("server: Site.Fallback is nil")
	}
	if s.Telemetry == nil {
		panic("server: Site.Telemetry is nil; use telemetry.Disabled for a site without telemetry")
	}
	p := s.Telemetry
	root := http.NewServeMux()
	root.Handle("GET /healthz", health(s.Version))
	root.Handle(telemetry.Path, p.Beacon())
	root.Handle("/", p.Handler(security(s.CSP, recovery(p.Logger(), s.Fallback, http.NewCrossOriginProtection().Handler(mux)))))
	return root
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

// health serves an uncached health response; mount it outside the middleware.
func health(version string) http.Handler {
	body := "ok " + version
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		io.WriteString(w, body)
	})
}
