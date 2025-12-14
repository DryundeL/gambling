package handlers

import (
	"encoding/json"
	"gambling/internal/application/use_case/spin"
	"log/slog"
	"net/http"
	"strconv"
)

// SpinHandler обрабатывает HTTP запросы для игры на спинах
type SpinHandler struct {
	spinUseCase *spin.SpinUseCase
	logger      *slog.Logger
}

// NewSpinHandler создает новый экземпляр SpinHandler
func NewSpinHandler(spinUseCase *spin.SpinUseCase, logger *slog.Logger) *SpinHandler {
	return &SpinHandler{
		spinUseCase: spinUseCase,
		logger:      logger,
	}
}

// SpinRequest представляет запрос на спин
type SpinRequest struct {
	BetAmount float64 `json:"bet_amount"`
}

// SpinResponse представляет ответ на спин
type SpinResponse struct {
	Reel1     int     `json:"reel1"`
	Reel2     int     `json:"reel2"`
	Reel3     int     `json:"reel3"`
	IsWin     bool    `json:"is_win"`
	WinAmount float64 `json:"win_amount"`
	Balance   float64 `json:"balance"`
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

	// Преобразуем HTTP запрос в команду use case
	cmd := spin.SpinCommand{
		UserID:    uint(userID),
		BetAmount: req.BetAmount,
	}

	// Выполняем use case
	result, err := h.spinUseCase.Execute(cmd)
	if err != nil {
		h.logger.Error("failed to spin", "error", err)
		if err.Error() == "неверная сумма" {
			http.Error(w, "Неверная сумма ставки", http.StatusBadRequest)
			return
		}
		if err.Error() == "недостаточно средств" {
			http.Error(w, "Недостаточно средств", http.StatusBadRequest)
			return
		}
		http.Error(w, "Ошибка при выполнении спина", http.StatusInternalServerError)
		return
	}

	response := SpinResponse{
		Reel1:     result.Reel1,
		Reel2:     result.Reel2,
		Reel3:     result.Reel3,
		IsWin:     result.IsWin,
		WinAmount: result.WinAmount,
		Balance:   result.Balance,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

