package app

import (
	"gambling/internal/application/use_case/auth"
	"gambling/internal/application/use_case/balance"
	"gambling/internal/application/use_case/spin"
	"gambling/internal/config"
	spinDomain "gambling/internal/domain/spin"
	"gambling/internal/infrastructure/database/pgsql"
	consoleInterface "gambling/internal/interfaces/console"
	"gambling/internal/infrastructure/repository"
	"log/slog"
)

// NewConsoleApp создает консольное приложение с использованием DDD архитектуры
func NewConsoleApp(cfg *config.Config, log *slog.Logger) *consoleInterface.Console {
	storage := pgsql.New(cfg)

	// Инициализация инфраструктуры (репозитории)
	userRepo := repository.NewUserRepository(storage.DB)
	transactionRepo := repository.NewTransactionRepository(storage.DB)
	spinRepo := repository.NewSpinRepository(storage.DB)

	// Инициализация доменного слоя
	spinDomainService := spinDomain.NewService()

	// Инициализация application слоя (use cases)
	registerUseCase := auth.NewRegisterUseCase(userRepo)
	loginUseCase := auth.NewLoginUseCase(userRepo)
	depositUseCase := balance.NewDepositUseCase(userRepo, transactionRepo)
	spinUC := spin.NewSpinUseCase(userRepo, transactionRepo, spinRepo, spinDomainService)

	// Создаем консольный интерфейс
	return consoleInterface.NewConsole(registerUseCase, loginUseCase, depositUseCase, spinUC)
}
