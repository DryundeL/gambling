package model

import (
	"time"

	"gorm.io/gorm"
)

// SpinResult представляет результат игры на спинах
type SpinResult struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	UserID        uint           `gorm:"not null;index" json:"user_id"`
	BetAmount     float64        `gorm:"not null;type:decimal(15,2)" json:"bet_amount"`
	WinAmount     float64        `gorm:"not null;type:decimal(15,2)" json:"win_amount"`
	Reel1         int            `gorm:"not null" json:"reel1"` // Символ на первом барабане (0-9)
	Reel2         int            `gorm:"not null" json:"reel2"` // Символ на втором барабане (0-9)
	Reel3         int            `gorm:"not null" json:"reel3"` // Символ на третьем барабане (0-9)
	IsWin         bool           `gorm:"not null" json:"is_win"`
	CreatedAt     time.Time      `json:"created_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName возвращает имя таблицы для модели SpinResult
func (SpinResult) TableName() string {
	return "spin_results"
}

