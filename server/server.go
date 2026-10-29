// Package server runs HTTP sites behind Cloudflare with explicit middleware.
// Given an application mux and a fallback handler, compose the chain as follows:
//
//	h := http.NewCrossOriginProtection().Handler(mux)
//	h = server.Recover(log, fallback, h)
//	h = server.Security(csp, h)
//	h = server.Log(log, h)
//	root := http.NewServeMux()
//	root.Handle("GET /healthz", server.Health(version))
//	root.Handle("/", h)
//	err := server.Run(ctx, server.New(":8080", root, log))
//
// Keep the same *http.Request between Log and mux so Log sees the route pattern.
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
)

var drainTimeout = 10 * time.Second

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

// Run binds synchronously, then serves until failure or ctx cancellation.
// On cancellation it drains for up to 10 seconds, then closes on drain failure.
// Close does not wait for handlers, so they may outlive Run after the deadline.
// Closing their resources after Run returns can race those handlers.
func Run(ctx context.Context, srv *http.Server, after ...func(context.Context) error) error {
	ln, err := net.Listen("tcp", cmp.Or(srv.Addr, ":http"))
	if err != nil {
		return err
	}
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ln) }()

	select {
	case err = <-served:
	case <-ctx.Done():
	}

	drain, cancel := context.WithTimeout(context.WithoutCancel(ctx), drainTimeout)
	defer cancel()
	if err == nil {
		if err = srv.Shutdown(drain); err != nil {
			err = errors.Join(err, srv.Close())
		}
		if served := <-served; !errors.Is(served, http.ErrServerClosed) {
			err = errors.Join(err, served)
		}
	} else if errors.Is(err, http.ErrServerClosed) {
		err = nil
	}
	for _, f := range after {
		err = errors.Join(err, f(drain))
	}
	return err
}

// Health serves an uncached health response; mount it outside the middleware.
func Health(version string) http.Handler {
	body := "ok " + version
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		io.WriteString(w, body)
	})
}
