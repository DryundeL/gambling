package router

import (
	"gambling/internal/database/pgsql"
	"net/http"

	mvLog "gambling/internal/http/middleware/logger"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New создаёт новый Router с подключенными хэндлерами.
// Параметры:
// - storage: экземпляр вашего pgsql хранилища
// - logger: ваш логгер для логирования запросов и ошибок
func New(storage *pgsql.Storage, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(mvLog.New(logger))

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte("OK")); err != nil {
			logger.Error("failed to write health check response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	return r
}
