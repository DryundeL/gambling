package router

import (
	"gambling/internal/application/use_case/auth"
	"gambling/internal/application/use_case/balance"
	spinUseCase "gambling/internal/application/use_case/spin"
	"gambling/internal/domain/spin"
	"gambling/internal/infrastructure/database/pgsql"
	"gambling/internal/interfaces/http/handlers"
	"gambling/internal/infrastructure/repository"
	"net/http"

	mvLog "gambling/internal/interfaces/http/middleware/logger"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// New создаёт новый Router с подключенными хэндлерами
// Здесь происходит композиция всех слоев DDD архитектуры
func New(storage *pgsql.Storage, logger *slog.Logger) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(mvLog.New(logger))
	r.Use(middleware.Logger)

	// ============================================
	// ИНИЦИАЛИЗАЦИЯ ИНФРАСТРУКТУРЫ (Infrastructure Layer)
	// ============================================
	// Создаем репозитории - это адаптеры для работы с БД
	userRepo := repository.NewUserRepository(storage.DB)
	transactionRepo := repository.NewTransactionRepository(storage.DB)
	spinRepo := repository.NewSpinRepository(storage.DB)

	// ============================================
	// ИНИЦИАЛИЗАЦИЯ ДОМЕННОГО СЛОЯ (Domain Layer)
	// ============================================
	// Создаем доменный сервис для логики игры
	spinDomainService := spin.NewService()

	// ============================================
	// ИНИЦИАЛИЗАЦИЯ APPLICATION СЛОЯ (Application Layer)
	// ============================================
	// Создаем use cases - это бизнес-операции приложения
	registerUseCase := auth.NewRegisterUseCase(userRepo)
	loginUseCase := auth.NewLoginUseCase(userRepo)
	depositUseCase := balance.NewDepositUseCase(userRepo, transactionRepo)
	spinUC := spinUseCase.NewSpinUseCase(userRepo, transactionRepo, spinRepo, spinDomainService)

	// ============================================
	// ИНИЦИАЛИЗАЦИЯ INTERFACES СЛОЯ (Interfaces Layer)
	// ============================================
	// Создаем HTTP handlers - это адаптеры для HTTP протокола
	authHandler := handlers.NewAuthHandler(registerUseCase, loginUseCase, logger)
	balanceHandler := handlers.NewBalanceHandler(depositUseCase, logger)
	spinHandler := handlers.NewSpinHandler(spinUC, logger)

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

