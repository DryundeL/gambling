package user

import (
	"time"

	"gorm.io/gorm"
)

// User представляет доменную сущность пользователя
// В DDD это Entity - объект с уникальной идентичностью
type User struct {
	ID           uint
	Username     string
	Email        string
	PasswordHash string
	Balance      float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt
}

// NewUser создает нового пользователя с начальным балансом 0
func NewUser(username, email, passwordHash string) *User {
	return &User{
		Username:     username,
		Email:        email,
		PasswordHash: passwordHash,
		Balance:      0,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// Deposit пополняет баланс пользователя
// Это доменная логика, которая инкапсулирована в сущности
func (u *User) Deposit(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	u.Balance += amount
	u.UpdatedAt = time.Now()
	return nil
}

// Withdraw списывает средства с баланса пользователя
func (u *User) Withdraw(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	if u.Balance < amount {
		return ErrInsufficientFunds
	}
	u.Balance -= amount
	u.UpdatedAt = time.Now()
	return nil
}

// AddWin добавляет выигрыш на баланс
func (u *User) AddWin(amount float64) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}
	u.Balance += amount
	u.UpdatedAt = time.Now()
	return nil
}

