package app

import (
	"context"
	"errors"
	"gambling/internal/config"
	"gambling/internal/database/pgsql"
	"gambling/internal/http/router"
	"log/slog"
	"net/http"
)

type App struct {
	cfg     *config.Config
	log     *slog.Logger
	storage *pgsql.Storage
	port    string
	routes  http.Handler
	server  *http.Server
}

func NewApp(cfg *config.Config, log *slog.Logger) *App {
	storage := pgsql.New(cfg)
	routes := router.New(storage, log)

	server := &http.Server{
		Addr:    cfg.AppUrl + ":" + cfg.AppPort,
		Handler: routes,
	}

	return &App{
		cfg:     cfg,
		log:     log,
		storage: storage,
		port:    cfg.AppPort,
		routes:  routes,
		server:  server,
	}
}

func (a *App) MustRun() {
	const op = "app.MustRun"

	log := a.log.With(slog.String("operation", op), slog.String("port", a.port))

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", slog.Any("error", err))
		}
	}()
}

func (a *App) Shutdown(ctx context.Context) error {
	const op = "app.Shutdown"

	log := a.log.With(slog.String("operation", op))
	log.Info("shutting down server...")

	if a.server == nil {
		return errors.New("server is not initialized")
	}

	return a.server.Shutdown(ctx)
}
