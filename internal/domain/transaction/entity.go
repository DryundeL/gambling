package transaction

import "time"

// Type определяет тип транзакции
type Type string

const (
	TypeDeposit Type = "deposit" // Пополнение
	TypeSpin    Type = "spin"    // Ставка в игре
	TypeWin     Type = "win"     // Выигрыш
)

// Transaction представляет доменную сущность транзакции
// Транзакция - это запись о финансовой операции пользователя
type Transaction struct {
	ID            uint
	UserID        uint
	Type          Type
	Amount        float64
	BalanceBefore float64
	BalanceAfter  float64
	Description   string
	CreatedAt     time.Time
}

// NewTransaction создает новую транзакцию
func NewTransaction(userID uint, txType Type, amount, balanceBefore, balanceAfter float64, description string) *Transaction {
	return &Transaction{
		UserID:        userID,
		Type:          txType,
		Amount:        amount,
		BalanceBefore: balanceBefore,
		BalanceAfter:  balanceAfter,
		Description:   description,
		CreatedAt:     time.Now(),
	}
}

