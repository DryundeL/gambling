package model

import (
	"time"

	"gorm.io/gorm"
)

// TransactionType определяет тип транзакции
type TransactionType string

const (
	TransactionTypeDeposit TransactionType = "deposit" // Пополнение
	TransactionTypeSpin    TransactionType = "spin"    // Игра на спинах
	TransactionTypeWin     TransactionType = "win"      // Выигрыш
)

// Transaction представляет транзакцию пользователя
type Transaction struct {
	ID            uint            `gorm:"primaryKey" json:"id"`
	UserID        uint            `gorm:"not null;index" json:"user_id"`
	Type          TransactionType `gorm:"not null;type:varchar(20)" json:"type"`
	Amount        float64         `gorm:"not null;type:decimal(15,2)" json:"amount"`
	BalanceBefore float64         `gorm:"not null;type:decimal(15,2)" json:"balance_before"`
	BalanceAfter  float64         `gorm:"not null;type:decimal(15,2)" json:"balance_after"`
	Description   string          `gorm:"size:255" json:"description"`
	CreatedAt     time.Time       `json:"created_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"index" json:"-"`
}

// TableName возвращает имя таблицы для модели Transaction
func (Transaction) TableName() string {
	return "transactions"
}

