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

func New(addr string, h http.Handler, log *slog.Logger) *http.Server {
	return &http.Server{
		Addr:              addr,
		Handler:           h,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       2 * time.Minute,
		ErrorLog:          slog.NewLogLogger(log.Handler(), slog.LevelWarn),
	}
}

func Run(ctx context.Context, srv *http.Server) error {
	ln, err := net.Listen("tcp", cmp.Or(srv.Addr, ":http"))
	if err != nil {
		return err
	}
	served := make(chan error, 1)
	go func() { served <- srv.Serve(ln) }()

	select {
	case err := <-served:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
	}

	drain, cancel := context.WithTimeout(context.WithoutCancel(ctx), drainTimeout)
	defer cancel()
	err = srv.Shutdown(drain)
	if err != nil {
		err = errors.Join(err, srv.Close())
	}
	if served := <-served; !errors.Is(served, http.ErrServerClosed) {
		err = errors.Join(err, served)
	}
	return err
}

func Health(version string) http.Handler {
	body := "ok " + version
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store")
		io.WriteString(w, body)
	})
}
