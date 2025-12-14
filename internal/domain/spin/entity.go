package spin

import "time"

// Result представляет доменную сущность результата спина
// Это запись о результате игры пользователя
type Result struct {
	ID        uint
	UserID    uint
	BetAmount float64
	WinAmount float64
	Reel1     int // Символ на первом барабане (0-9)
	Reel2     int // Символ на втором барабане (0-9)
	Reel3     int // Символ на третьем барабане (0-9)
	IsWin     bool
	CreatedAt time.Time
}

// NewResult создает новый результат спина
func NewResult(userID uint, betAmount, winAmount float64, reel1, reel2, reel3 int) *Result {
	return &Result{
		UserID:    userID,
		BetAmount: betAmount,
		WinAmount: winAmount,
		Reel1:     reel1,
		Reel2:     reel2,
		Reel3:     reel3,
		IsWin:     winAmount > 0,
		CreatedAt: time.Now(),
	}
}

