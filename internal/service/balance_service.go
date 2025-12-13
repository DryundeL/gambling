package service

import (
	"errors"
	"gambling/internal/model"
	"gambling/internal/repository"
)

var (
	ErrInvalidAmount = errors.New("неверная сумма")
	ErrInsufficientFunds = errors.New("недостаточно средств")
)

// BalanceService предоставляет методы для работы с балансом пользователя
type BalanceService struct {
	userRepo        *repository.UserRepository
	transactionRepo *repository.TransactionRepository
}

// NewBalanceService создает новый экземпляр BalanceService
func NewBalanceService(userRepo *repository.UserRepository, transactionRepo *repository.TransactionRepository) *BalanceService {
	return &BalanceService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
	}
}

// Deposit пополняет баланс пользователя
func (s *BalanceService) Deposit(userID uint, amount float64) (*model.User, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	balanceBefore := user.Balance
	balanceAfter := balanceBefore + amount

	// Обновляем баланс
	if err := s.userRepo.UpdateBalance(userID, balanceAfter); err != nil {
		return nil, err
	}

	// Создаем транзакцию
	transaction := &model.Transaction{
		UserID:        userID,
		Type:          model.TransactionTypeDeposit,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   "Пополнение баланса",
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		// Откатываем баланс в случае ошибки
		_ = s.userRepo.UpdateBalance(userID, balanceBefore)
		return nil, err
	}

	// Получаем обновленного пользователя
	user, err = s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// Withdraw списывает средства с баланса пользователя
func (s *BalanceService) Withdraw(userID uint, amount float64) (*model.User, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	if user.Balance < amount {
		return nil, ErrInsufficientFunds
	}

	balanceBefore := user.Balance
	balanceAfter := balanceBefore - amount

	// Обновляем баланс
	if err := s.userRepo.UpdateBalance(userID, balanceAfter); err != nil {
		return nil, err
	}

	// Создаем транзакцию
	transaction := &model.Transaction{
		UserID:        userID,
		Type:          model.TransactionTypeSpin,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   "Ставка в игре",
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		// Откатываем баланс в случае ошибки
		_ = s.userRepo.UpdateBalance(userID, balanceBefore)
		return nil, err
	}

	// Получаем обновленного пользователя
	user, err = s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

// AddWin добавляет выигрыш на баланс пользователя
func (s *BalanceService) AddWin(userID uint, amount float64) (*model.User, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	user, err := s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	balanceBefore := user.Balance
	balanceAfter := balanceBefore + amount

	// Обновляем баланс
	if err := s.userRepo.UpdateBalance(userID, balanceAfter); err != nil {
		return nil, err
	}

	// Создаем транзакцию
	transaction := &model.Transaction{
		UserID:        userID,
		Type:          model.TransactionTypeWin,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   "Выигрыш в игре",
	}

	if err := s.transactionRepo.Create(transaction); err != nil {
		// Откатываем баланс в случае ошибки
		_ = s.userRepo.UpdateBalance(userID, balanceBefore)
		return nil, err
	}

	// Получаем обновленного пользователя
	user, err = s.userRepo.GetByID(userID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

