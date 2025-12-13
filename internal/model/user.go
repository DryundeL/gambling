package model

import (
	"time"

	"gorm.io/gorm"
)

// User представляет модель пользователя в системе казино
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"uniqueIndex;not null;size:50" json:"username"`
	Email        string         `gorm:"uniqueIndex;not null;size:100" json:"email"`
	PasswordHash string         `gorm:"not null;size:255" json:"-"`
	Balance      float64        `gorm:"not null;default:0;type:decimal(15,2)" json:"balance"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName возвращает имя таблицы для модели User
func (User) TableName() string {
	return "users"
}
