package handlers

import (
	"encoding/json"
	"gambling/internal/service"
	"log/slog"
	"net/http"
	"strconv"
)

// SpinHandler обрабатывает HTTP запросы для игры на спинах
type SpinHandler struct {
	spinService *service.SpinService
	logger      *slog.Logger
}

// NewSpinHandler создает новый экземпляр SpinHandler
func NewSpinHandler(spinService *service.SpinService, logger *slog.Logger) *SpinHandler {
	return &SpinHandler{
		spinService: spinService,
		logger:      logger,
	}
}

// SpinRequest представляет запрос на спин
type SpinRequest struct {
	BetAmount float64 `json:"bet_amount"`
}

// Spin обрабатывает запрос на выполнение спина
func (h *SpinHandler) Spin(w http.ResponseWriter, r *http.Request) {
	// Получаем userID из параметров запроса (в реальном приложении из JWT токена)
	userIDStr := r.URL.Query().Get("user_id")
	if userIDStr == "" {
		http.Error(w, "user_id обязателен", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseUint(userIDStr, 10, 32)
	if err != nil {
		http.Error(w, "Неверный формат user_id", http.StatusBadRequest)
		return
	}

	var req SpinRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	result, err := h.spinService.Spin(uint(userID), req.BetAmount)
	if err != nil {
		h.logger.Error("failed to spin", "error", err)
		if err == service.ErrInvalidBetAmount {
			http.Error(w, "Неверная сумма ставки", http.StatusBadRequest)
			return
		}
		if err == service.ErrInsufficientFunds {
			http.Error(w, "Недостаточно средств", http.StatusBadRequest)
			return
		}
		http.Error(w, "Ошибка при выполнении спина", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(result); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

