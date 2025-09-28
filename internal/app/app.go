package app

import (
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
}

func NewApp(cfg *config.Config, log *slog.Logger) *App {
	storage := pgsql.New(cfg)
	routes := router.New(storage, log)

	return &App{
		cfg:     cfg,
		log:     log,
		storage: storage,
		port:    cfg.AppPort,
		routes:  routes,
	}
}

func (a *App) MustRun() {

}

func (a *App) Run() error {
	const op = "gprcapp.Run"

	log := a.log.With(slog.String("operation", op), slog.String("port", a.port))

	server := http.Server{
		Addr:    a.cfg.AppUrl + ":" + a.cfg.AppPort,
		Handler: a.routes,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("failed to start server", slog.Any("error", err))
		}
	}()

	return nil
}
