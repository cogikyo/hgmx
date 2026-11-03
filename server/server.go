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

	"connectrpc.com/connect/v2"

	"github.com/cogikyo/hgmx/telemetry"
)

var drainTimeout = 5 * time.Second

const stopTimeout = 8 * time.Second

type Site struct {
	Version string
	CSP     CSP
	// Allow grants (self) for named default Permissions-Policy features.
	Allow []string
	// InternalError is fixed panic-page HTML; empty uses hgmx's default.
	// Recovery serves it as 500 without rendering templates.
	InternalError string
	Telemetry     *telemetry.Pipeline
	// RPC identifies mounted services for RPC recovery and unknown-service faults.
	// Set it after registration; it does not mount handlers.
	RPC *connect.Server
}

// Handler wraps mux with telemetry, panic recovery and cross-origin protection.
// The mux type rejects outer wrappers that can hide r.Pattern from telemetry.
// StripPrefix, nested muxes, and route wrappers can still coarsen http.route.
// Security headers cover all routes; health and enabled beacons bypass mux.
// It panics on nil Telemetry, an invalid CSP, or unknown Allow names.
// Use [telemetry.Disabled] to opt out of telemetry.
func (s Site) Handler(mux *http.ServeMux) http.Handler {
	if s.Telemetry == nil {
		panic("server: Site.Telemetry is nil; use telemetry.Disabled for a site without telemetry")
	}
	p := s.Telemetry
	root := http.NewServeMux()
	root.Handle("GET /healthz", health(s.Version))
	if b := p.Beacon(); b.Head() != "" {
		root.Handle(telemetry.Path, b)
	}
	fallback, routes := internalError(s.InternalError), http.Handler(mux)
	if s.RPC != nil {
		rpc := newRPC(s.RPC)
		fallback, routes = rpc.fallback(fallback), rpc.route(mux)
	}
	app := p.Handler(recovery(p.Logger(), fallback, http.NewCrossOriginProtection().Handler(routes)))
	return security(s.CSP.Header(), s.Allow, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if h, pattern := root.Handler(r); pattern != "" {
			h.ServeHTTP(w, r)
			return
		}
		app.ServeHTTP(w, r)
	}))
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
// After a successful bind, it drains for up to 5 seconds before the hooks.
// On drain failure it closes connections; handlers may still outlive Run.
// Closing their resources after Run returns can race those handlers.
// The after hooks share a deadline 8 seconds after shutdown starts.
// They run even if binding fails, with the full 8 seconds.
// Run returns all errors joined.
func Run(ctx context.Context, srv *http.Server, after ...func(context.Context) error) error {
	stop, err := serve(ctx, srv)
	hooks, cancel := context.WithDeadline(context.WithoutCancel(ctx), stop.Add(stopTimeout))
	defer cancel()
	for _, f := range after {
		err = errors.Join(err, f(hooks))
	}
	return err
}

func serve(ctx context.Context, srv *http.Server) (stop time.Time, err error) {
	ln, err := net.Listen("tcp", cmp.Or(srv.Addr, ":http"))
	if err != nil {
		return time.Now(), err
	}
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ln) }()

	var serr error
	select {
	case serr = <-served:
	case <-ctx.Done():
	}

	stop = time.Now()
	drain, cancel := context.WithDeadline(context.WithoutCancel(ctx), stop.Add(drainTimeout))
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
	return stop, err
}

func health(version string) http.Handler {
	body := "ok " + version
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		io.WriteString(w, body)
	})
}
