package balance

import (
	"gambling/internal/domain/transaction"
	"gambling/internal/domain/user"
)

// DepositUseCase представляет use case для пополнения баланса
type DepositUseCase struct {
	userRepo        user.Repository
	transactionRepo transaction.Repository
}

// NewDepositUseCase создает новый use case для пополнения баланса
func NewDepositUseCase(userRepo user.Repository, transactionRepo transaction.Repository) *DepositUseCase {
	return &DepositUseCase{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

// DepositCommand представляет команду для пополнения баланса
type DepositCommand struct {
	UserID uint
	Amount float64
}

// DepositResult представляет результат пополнения баланса
type DepositResult struct {
	Balance float64
}

// Execute выполняет пополнение баланса пользователя
func (uc *DepositUseCase) Execute(cmd DepositCommand) (*DepositResult, error) {
	// Получаем пользователя
	u, err := uc.userRepo.GetByID(cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Выполняем доменную операцию пополнения
	balanceBefore := u.Balance
	if err := u.Deposit(cmd.Amount); err != nil {
		return nil, err
	}

	// Сохраняем обновленный баланс
	if err := uc.userRepo.UpdateBalance(cmd.UserID, u.Balance); err != nil {
		return nil, err
	}

	// Создаем транзакцию
	tx := transaction.NewTransaction(
		cmd.UserID,
		transaction.TypeDeposit,
		cmd.Amount,
		balanceBefore,
		u.Balance,
		"Пополнение баланса",
	)

	if err := uc.transactionRepo.Create(tx); err != nil {
		// Откатываем баланс в случае ошибки
		_ = uc.userRepo.UpdateBalance(cmd.UserID, balanceBefore)
		return nil, err
	}

	return &DepositResult{
		Balance: u.Balance,
	}, nil
}

