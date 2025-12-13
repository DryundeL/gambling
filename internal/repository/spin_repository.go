package repository

import (
	"gambling/internal/model"

	"gorm.io/gorm"
)

// SpinRepository предоставляет методы для работы с результатами спинов
type SpinRepository struct {
	db *gorm.DB
}

// NewSpinRepository создает новый экземпляр SpinRepository
func NewSpinRepository(db *gorm.DB) *SpinRepository {
	return &SpinRepository{db: db}
}

// Create создает новый результат спина
func (r *SpinRepository) Create(spinResult *model.SpinResult) error {
	return r.db.Create(spinResult).Error
}

// GetByUserID возвращает историю спинов пользователя
func (r *SpinRepository) GetByUserID(userID uint, limit int) ([]model.SpinResult, error) {
	var spins []model.SpinResult
	query := r.db.Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&spins).Error; err != nil {
		return nil, err
	}
	return spins, nil
}

