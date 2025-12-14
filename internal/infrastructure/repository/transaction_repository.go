package repository

import (
	"gambling/internal/domain/transaction"
	"time"

	"gorm.io/gorm"
)

// TransactionRepository реализует интерфейс transaction.Repository
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository создает новый репозиторий транзакций
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create создает новую транзакцию
func (r *TransactionRepository) Create(tx *transaction.Transaction) error {
	dbTx := toDBTransaction(tx)
	if err := r.db.Create(dbTx).Error; err != nil {
		return err
	}
	tx.ID = dbTx.ID
	tx.CreatedAt = dbTx.CreatedAt
	return nil
}

// GetByUserID возвращает все транзакции пользователя
func (r *TransactionRepository) GetByUserID(userID uint, limit int) ([]*transaction.Transaction, error) {
	var dbTxs []DBTransaction
	query := r.db.Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&dbTxs).Error; err != nil {
		return nil, err
	}

	result := make([]*transaction.Transaction, len(dbTxs))
	for i, dbTx := range dbTxs {
		result[i] = toDomainTransaction(&dbTx)
	}
	return result, nil
}

// DBTransaction представляет модель БД для транзакции
type DBTransaction struct {
	ID            uint            `gorm:"primaryKey"`
	UserID        uint            `gorm:"not null;index"`
	Type          string          `gorm:"not null;type:varchar(20)"`
	Amount        float64         `gorm:"not null;type:decimal(15,2)"`
	BalanceBefore float64         `gorm:"not null;type:decimal(15,2)"`
	BalanceAfter  float64         `gorm:"not null;type:decimal(15,2)"`
	Description   string          `gorm:"size:255"`
	CreatedAt     time.Time       `gorm:"autoCreateTime"`
	DeletedAt     gorm.DeletedAt  `gorm:"index"`
}

func (DBTransaction) TableName() string {
	return "transactions"
}

func toDBTransaction(tx *transaction.Transaction) *DBTransaction {
	return &DBTransaction{
		ID:            tx.ID,
		UserID:        tx.UserID,
		Type:          string(tx.Type),
		Amount:        tx.Amount,
		BalanceBefore: tx.BalanceBefore,
		BalanceAfter:  tx.BalanceAfter,
		Description:   tx.Description,
		CreatedAt:     tx.CreatedAt,
	}
}

func toDomainTransaction(dbTx *DBTransaction) *transaction.Transaction {
	return &transaction.Transaction{
		ID:            dbTx.ID,
		UserID:        dbTx.UserID,
		Type:          transaction.Type(dbTx.Type),
		Amount:        dbTx.Amount,
		BalanceBefore: dbTx.BalanceBefore,
		BalanceAfter:  dbTx.BalanceAfter,
		Description:   dbTx.Description,
		CreatedAt:     dbTx.CreatedAt,
	}
}

