package service

import (
	"errors"
	"gambling/internal/model"
	"gambling/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("неверные учетные данные")
	ErrUserExists         = errors.New("пользователь уже существует")
)

// AuthService предоставляет методы для аутентификации и регистрации
type AuthService struct {
	userRepo *repository.UserRepository
}

// NewAuthService создает новый экземпляр AuthService
func NewAuthService(userRepo *repository.UserRepository) *AuthService {
	return &AuthService{userRepo: userRepo}
}

// Register регистрирует нового пользователя
func (s *AuthService) Register(username, email, password string) (*model.User, error) {
	// Проверяем, существует ли пользователь с таким username
	if _, err := s.userRepo.GetByUsername(username); err == nil {
		return nil, ErrUserExists
	}

	// Проверяем, существует ли пользователь с таким email
	if _, err := s.userRepo.GetByEmail(email); err == nil {
		return nil, ErrUserExists
	}

	// Хешируем пароль
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Создаем пользователя
	user := &model.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		Balance:      0,
	}

	if err := s.userRepo.Create(user); err != nil {
		if err == repository.ErrUserAlreadyExists {
			return nil, ErrUserExists
		}
		return nil, err
	}

	return user, nil
}

// Login проверяет учетные данные и возвращает пользователя
func (s *AuthService) Login(username, password string) (*model.User, error) {
	user, err := s.userRepo.GetByUsername(username)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Проверяем пароль
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

