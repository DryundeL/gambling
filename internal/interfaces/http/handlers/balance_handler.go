package handlers

import (
	"encoding/json"
	"gambling/internal/application/use_case/balance"
	"log/slog"
	"net/http"
	"strconv"
)

// BalanceHandler обрабатывает HTTP запросы для работы с балансом
type BalanceHandler struct {
	depositUseCase *balance.DepositUseCase
	logger         *slog.Logger
}

// NewBalanceHandler создает новый экземпляр BalanceHandler
func NewBalanceHandler(depositUseCase *balance.DepositUseCase, logger *slog.Logger) *BalanceHandler {
	return &BalanceHandler{
		depositUseCase: depositUseCase,
		logger:         logger,
	}
}

// DepositRequest представляет запрос на пополнение баланса
type DepositRequest struct {
	Amount float64 `json:"amount"`
}

// DepositResponse представляет ответ на пополнение баланса
type DepositResponse struct {
	Balance float64 `json:"balance"`
}

// Deposit обрабатывает запрос на пополнение баланса
func (h *BalanceHandler) Deposit(w http.ResponseWriter, r *http.Request) {
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

	var req DepositRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("failed to decode request", "error", err)
		http.Error(w, "Неверный формат запроса", http.StatusBadRequest)
		return
	}

	// Преобразуем HTTP запрос в команду use case
	cmd := balance.DepositCommand{
		UserID: uint(userID),
		Amount: req.Amount,
	}

	// Выполняем use case
	result, err := h.depositUseCase.Execute(cmd)
	if err != nil {
		h.logger.Error("failed to deposit", "error", err)
		if err.Error() == "неверная сумма" {
			http.Error(w, "Неверная сумма", http.StatusBadRequest)
			return
		}
		http.Error(w, "Ошибка при пополнении баланса", http.StatusInternalServerError)
		return
	}

	response := DepositResponse{
		Balance: result.Balance,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

