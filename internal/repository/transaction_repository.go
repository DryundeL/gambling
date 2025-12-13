package repository

import (
	"gambling/internal/model"

	"gorm.io/gorm"
)

// TransactionRepository предоставляет методы для работы с транзакциями
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository создает новый экземпляр TransactionRepository
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create создает новую транзакцию
func (r *TransactionRepository) Create(transaction *model.Transaction) error {
	return r.db.Create(transaction).Error
}

// GetByUserID возвращает все транзакции пользователя
func (r *TransactionRepository) GetByUserID(userID uint, limit int) ([]model.Transaction, error) {
	var transactions []model.Transaction
	query := r.db.Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&transactions).Error; err != nil {
		return nil, err
	}
	return transactions, nil
}

