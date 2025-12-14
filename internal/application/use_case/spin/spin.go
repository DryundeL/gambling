package spin

import (
	"gambling/internal/domain/spin"
	"gambling/internal/domain/transaction"
	"gambling/internal/domain/user"
)

// SpinUseCase представляет use case для выполнения спина
type SpinUseCase struct {
	userRepo        user.Repository
	transactionRepo transaction.Repository
	spinRepo        spin.Repository
	spinService     *spin.Service
}

// NewSpinUseCase создает новый use case для спинов
func NewSpinUseCase(
	userRepo user.Repository,
	transactionRepo transaction.Repository,
	spinRepo spin.Repository,
	spinService *spin.Service,
) *SpinUseCase {
	return &SpinUseCase{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		spinRepo:        spinRepo,
		spinService:     spinService,
	}
}

// SpinCommand представляет команду для выполнения спина
type SpinCommand struct {
	UserID    uint
	BetAmount float64
}

// SpinResult представляет результат спина
type SpinResult struct {
	Reel1     int
	Reel2     int
	Reel3     int
	IsWin     bool
	WinAmount float64
	Balance   float64
}

// Execute выполняет спин игры
func (uc *SpinUseCase) Execute(cmd SpinCommand) (*SpinResult, error) {
	if cmd.BetAmount <= 0 {
		return nil, user.ErrInvalidAmount
	}

	// Получаем пользователя
	u, err := uc.userRepo.GetByID(cmd.UserID)
	if err != nil {
		return nil, err
	}

	// Списываем ставку через доменную логику
	balanceBefore := u.Balance
	if err := u.Withdraw(cmd.BetAmount); err != nil {
		return nil, err
	}

	// Сохраняем обновленный баланс
	if err := uc.userRepo.UpdateBalance(cmd.UserID, u.Balance); err != nil {
		return nil, err
	}

	// Создаем транзакцию на списание
	tx := transaction.NewTransaction(
		cmd.UserID,
		transaction.TypeSpin,
		cmd.BetAmount,
		balanceBefore,
		u.Balance,
		"Ставка в игре",
	)

	if err := uc.transactionRepo.Create(tx); err != nil {
		// Откатываем баланс в случае ошибки
		_ = uc.userRepo.UpdateBalance(cmd.UserID, balanceBefore)
		return nil, err
	}

	// Генерируем символы на барабанах через доменный сервис
	reel1 := uc.spinService.GenerateSymbol()
	reel2 := uc.spinService.GenerateSymbol()
	reel3 := uc.spinService.GenerateSymbol()

	// Вычисляем выигрыш через доменный сервис
	winAmount := uc.spinService.CalculateWin(reel1, reel2, reel3, cmd.BetAmount)
	isWin := winAmount > 0

	// Если есть выигрыш, добавляем его на баланс
	if isWin {
		balanceBeforeWin := u.Balance
		if err := u.AddWin(winAmount); err != nil {
			return nil, err
		}

		if err := uc.userRepo.UpdateBalance(cmd.UserID, u.Balance); err != nil {
			return nil, err
		}

		// Создаем транзакцию на выигрыш
		winTx := transaction.NewTransaction(
			cmd.UserID,
			transaction.TypeWin,
			winAmount,
			balanceBeforeWin,
			u.Balance,
			"Выигрыш в игре",
		)

		if err := uc.transactionRepo.Create(winTx); err != nil {
			// Откатываем баланс в случае ошибки
			_ = uc.userRepo.UpdateBalance(cmd.UserID, balanceBeforeWin)
			return nil, err
		}
	}

	// Сохраняем результат спина
	spinResult := spin.NewResult(cmd.UserID, cmd.BetAmount, winAmount, reel1, reel2, reel3)
	if err := uc.spinRepo.Create(spinResult); err != nil {
		return nil, err
	}

	// Получаем финальный баланс
	finalUser, err := uc.userRepo.GetByID(cmd.UserID)
	if err != nil {
		return nil, err
	}

	return &SpinResult{
		Reel1:     reel1,
		Reel2:     reel2,
		Reel3:     reel3,
		IsWin:     isWin,
		WinAmount: winAmount,
		Balance:   finalUser.Balance,
	}, nil
}

