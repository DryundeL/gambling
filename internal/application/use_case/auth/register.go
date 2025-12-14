package auth

import (
	"gambling/internal/domain/user"

	"golang.org/x/crypto/bcrypt"
)

// RegisterUseCase представляет use case для регистрации пользователя
// Use Case - это конкретная бизнес-операция, которую может выполнить пользователь
type RegisterUseCase struct {
	userRepo user.Repository
}

// NewRegisterUseCase создает новый use case для регистрации
func NewRegisterUseCase(userRepo user.Repository) *RegisterUseCase {
	return &RegisterUseCase{
		userRepo: userRepo,
	}
}

// RegisterCommand представляет команду для регистрации
// Command - это DTO для входных данных use case
type RegisterCommand struct {
	Username string
	Email    string
	Password string
}

// RegisterResult представляет результат регистрации
// Result - это DTO для выходных данных use case
type RegisterResult struct {
	ID       uint
	Username string
	Email    string
	Balance  float64
}

// Execute выполняет регистрацию нового пользователя
func (uc *RegisterUseCase) Execute(cmd RegisterCommand) (*RegisterResult, error) {
	// Проверяем, существует ли пользователь с таким username
	if _, err := uc.userRepo.GetByUsername(cmd.Username); err == nil {
		return nil, user.ErrUserAlreadyExists
	}

	// Проверяем, существует ли пользователь с таким email
	if _, err := uc.userRepo.GetByEmail(cmd.Email); err == nil {
		return nil, user.ErrUserAlreadyExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(cmd.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем доменную сущность пользователя
	newUser := user.NewUser(cmd.Username, cmd.Email, string(hashedPassword))

	// Сохраняем через репозиторий
	if err := uc.userRepo.Create(newUser); err != nil {
		return nil, err
	}

	return &RegisterResult{
		ID:       newUser.ID,
		Username: newUser.Username,
		Email:    newUser.Email,
		Balance:  newUser.Balance,
	}, nil
}

