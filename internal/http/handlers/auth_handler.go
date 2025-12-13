package handlers

import (
	"encoding/json"
	"gambling/internal/service"
	"log/slog"
	"net/http"
)

// AuthHandler обрабатывает HTTP запросы для аутентификации
type AuthHandler struct {
	authService *service.AuthService
	logger      *slog.Logger
}

// NewAuthHandler создает новый экземпляр AuthHandler
func NewAuthHandler(authService *service.AuthService, logger *slog.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
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

	user, err := h.authService.Register(req.Username, req.Email, req.Password)
	if err != nil {
		h.logger.Error("failed to register user", "error", err)
		if err == service.ErrUserExists {
			http.Error(w, "Пользователь уже существует", http.StatusConflict)
			return
		}
		http.Error(w, "Ошибка при регистрации", http.StatusInternalServerError)
		return
	}

	response := RegisterResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Balance:  user.Balance,
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

	user, err := h.authService.Login(req.Username, req.Password)
	if err != nil {
		h.logger.Error("failed to login user", "error", err)
		if err == service.ErrInvalidCredentials {
			http.Error(w, "Неверные учетные данные", http.StatusUnauthorized)
			return
		}
		http.Error(w, "Ошибка при входе", http.StatusInternalServerError)
		return
	}

	response := LoginResponse{
		ID:       user.ID,
		Username: user.Username,
		Email:    user.Email,
		Balance:  user.Balance,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

