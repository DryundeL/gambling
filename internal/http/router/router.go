package router

import (
	"gambling/internal/database/pgsql"
	"gambling/internal/http/handlers"
	"gambling/internal/repository"
	"gambling/internal/service"
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
	r.Use(middleware.Logger)

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(storage.DB)
	transactionRepo := repository.NewTransactionRepository(storage.DB)
	spinRepo := repository.NewSpinRepository(storage.DB)

	// Инициализация сервисов
	authService := service.NewAuthService(userRepo)
	balanceService := service.NewBalanceService(userRepo, transactionRepo)
	spinService := service.NewSpinService(balanceService, spinRepo)

	// Инициализация хэндлеров
	authHandler := handlers.NewAuthHandler(authService, logger)
	balanceHandler := handlers.NewBalanceHandler(balanceService, logger)
	spinHandler := handlers.NewSpinHandler(spinService, logger)

	// Маршруты
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		if _, err := w.Write([]byte("OK")); err != nil {
			logger.Error("failed to write health check response", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}
	})

	// API маршруты
	r.Route("/api/v1", func(r chi.Router) {
		// Аутентификация
		r.Post("/register", authHandler.Register)
		r.Post("/login", authHandler.Login)

		// Баланс
		r.Post("/balance/deposit", balanceHandler.Deposit)

		// Игра
		r.Post("/spin", spinHandler.Spin)
	})

	return r
}
