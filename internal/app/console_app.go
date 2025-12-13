package app

import (
	"gambling/internal/config"
	"gambling/internal/console"
	"gambling/internal/database/pgsql"
	"gambling/internal/repository"
	"gambling/internal/service"
	"log/slog"
)

// NewConsoleApp создает консольное приложение
func NewConsoleApp(cfg *config.Config, log *slog.Logger) *console.Console {
	storage := pgsql.New(cfg)

	// Инициализация репозиториев
	userRepo := repository.NewUserRepository(storage.DB)
	transactionRepo := repository.NewTransactionRepository(storage.DB)
	spinRepo := repository.NewSpinRepository(storage.DB)

	// Инициализация сервисов
	authService := service.NewAuthService(userRepo)
	balanceService := service.NewBalanceService(userRepo, transactionRepo)
	spinService := service.NewSpinService(balanceService, spinRepo)

	// Создаем консольный интерфейс
	return console.NewConsole(authService, balanceService, spinService)
}
