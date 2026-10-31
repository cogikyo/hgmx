package routes

import (
	"context"
	"log/slog"
	"net/http"
)

func Home(w http.ResponseWriter, r *http.Request) {
	slog.Info("home")
	slog.InfoContext(r.Context(), "home")
}

func Sync(ctx context.Context, log *slog.Logger) {
	go func() {
		log.Warn("sync")
	}()
}

func Start() {
	slog.Info("start")
	ctx := context.Background()
	slog.Debug("ready", "ctx", ctx)
}
