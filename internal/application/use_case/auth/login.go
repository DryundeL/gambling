package auth

import (
	"gambling/internal/domain/user"

	"golang.org/x/crypto/bcrypt"
)

// LoginUseCase представляет use case для входа пользователя
type LoginUseCase struct {
	userRepo user.Repository
}

// NewLoginUseCase создает новый use case для входа
func NewLoginUseCase(userRepo user.Repository) *LoginUseCase {
	return &LoginUseCase{
		userRepo: userRepo,
	}
}

// LoginCommand представляет команду для входа
type LoginCommand struct {
	Username string
	Password string
}

// LoginResult представляет результат входа
type LoginResult struct {
	ID       uint
	Username string
	Email    string
	Balance  float64
}

// Execute выполняет вход пользователя
func (uc *LoginUseCase) Execute(cmd LoginCommand) (*LoginResult, error) {
	// Получаем пользователя по username
	u, err := uc.userRepo.GetByUsername(cmd.Username)
	if err != nil {
		return nil, user.ErrInvalidCredentials
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(cmd.Password)); err != nil {
		return nil, user.ErrInvalidCredentials
	}

	return &LoginResult{
		ID:       u.ID,
		Username: u.Username,
		Email:    u.Email,
		Balance:  u.Balance,
	}, nil
}
