package main

import (
	"context"
	"gambling/internal/app"
	"gambling/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.AppEnv)
	log.Info("starting app", slog.Any("cfg", cfg.AppUrl))

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	server := app.NewApp(cfg, log)
	server.MustRun()

	// Небольшая задержка для обеспечения запуска сервера
	time.Sleep(100 * time.Millisecond)
	log.Info("Server started", slog.String("addr", cfg.AppUrl+":"+cfg.AppPort))

	// Ожидаем сигнал для graceful shutdown
	sig := <-quit
	log.Info("shutting down server...", slog.Any("signal", sig))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.Any("error", err))
	} else {
		log.Info("server gracefully stopped")
	}
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}
