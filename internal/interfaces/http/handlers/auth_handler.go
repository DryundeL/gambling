package handlers

import (
	"encoding/json"
	"gambling/internal/application/use_case/auth"
	"log/slog"
	"net/http"
)

// AuthHandler обрабатывает HTTP запросы для аутентификации
// Это адаптер, который преобразует HTTP запросы в команды use case
type AuthHandler struct {
	registerUseCase *auth.RegisterUseCase
	loginUseCase    *auth.LoginUseCase
	logger          *slog.Logger
}

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(
	registerUseCase *auth.RegisterUseCase,
	loginUseCase *auth.LoginUseCase,
	logger *slog.Logger,
) *AuthHandler {
	return &AuthHandler{
		registerUseCase: registerUseCase,
		loginUseCase:    loginUseCase,
		logger:          logger,
	}
}

// RegisterRequest представляет запрос на регистрацию
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// RegisterResponse представляет ответ на регистрацию
type RegisterResponse struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Balance  float64 `json:"balance"`
}

// Register обрабатывает запрос на регистрацию
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "Все поля обязательны для заполнения", http.StatusBadRequest)
		return
	}

	// Преобразуем HTTP запрос в команду use case
	cmd := auth.RegisterCommand{
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	}

	// Выполняем use case
	result, err := h.registerUseCase.Execute(cmd)
	if err != nil {
		h.logger.Error("failed to register user", "error", err)
		if err.Error() == "пользователь уже существует" {
			http.Error(w, "Пользователь уже существует", http.StatusConflict)
			return
		}
		http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
		return
	}

	response := RegisterResponse{
		ID:       result.ID,
		Username: result.Username,
		Email:    result.Email,
		Balance:  result.Balance,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

// LoginRequest представляет запрос на вход
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse представляет ответ на вход
type LoginResponse struct {
	ID       uint    `json:"id"`
	Username string  `json:"username"`
	Email    string  `json:"email"`
	Balance  float64 `json:"balance"`
}

// Login обрабатывает запрос на вход
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	if req.Username == "" || req.Password == "" {
		http.Error(w, "Имя пользователя и пароль обязательны", http.StatusBadRequest)
		return
	}

	// Преобразуем HTTP запрос в команду use case
	cmd := auth.LoginCommand{
		Username: req.Username,
		Password: req.Password,
	}

	// Выполняем use case
	result, err := h.loginUseCase.Execute(cmd)
	if err != nil {
		h.logger.Error("failed to login user", "error", err)
		if err.Error() == "неверные учетные данные" {
			http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Ошибка при входе", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		ID:       result.ID,
		Username: result.Username,
		Email:    result.Email,
		Balance:  result.Balance,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

