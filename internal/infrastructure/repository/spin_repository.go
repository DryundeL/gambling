package repository

import (
	"gambling/internal/domain/spin"
	"time"

	"gorm.io/gorm"
)

// SpinRepository реализует интерфейс spin.Repository
type SpinRepository struct {
	db *gorm.DB
}

// NewSpinRepository создает новый репозиторий спинов
func NewSpinRepository(db *gorm.DB) *SpinRepository {
	return &SpinRepository{db: db}
}

// Create создает новый результат спина
func (r *SpinRepository) Create(result *spin.Result) error {
	dbResult := toDBSpinResult(result)
	if err := r.db.Create(dbResult).Error; err != nil {
		return err
	}
	result.ID = dbResult.ID
	result.CreatedAt = dbResult.CreatedAt
	return nil
}

// GetByUserID возвращает историю спинов пользователя
func (r *SpinRepository) GetByUserID(userID uint, limit int) ([]*spin.Result, error) {
	var dbResults []DBSpinResult
	query := r.db.Where("user_id = ?", userID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	if err := query.Find(&dbResults).Error; err != nil {
		return nil, err
	}

	result := make([]*spin.Result, len(dbResults))
	for i, dbResult := range dbResults {
		result[i] = toDomainSpinResult(&dbResult)
	}
	return result, nil
}

// DBSpinResult представляет модель БД для результата спина
type DBSpinResult struct {
	ID        uint           `gorm:"primaryKey"`
	UserID    uint           `gorm:"not null;index"`
	BetAmount float64        `gorm:"not null;type:decimal(15,2)"`
	WinAmount float64        `gorm:"not null;type:decimal(15,2)"`
	Reel1     int            `gorm:"not null"`
	Reel2     int            `gorm:"not null"`
	Reel3     int            `gorm:"not null"`
	IsWin     bool           `gorm:"not null"`
	CreatedAt time.Time      `gorm:"autoCreateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (DBSpinResult) TableName() string {
	return "spin_results"
}

func toDBSpinResult(result *spin.Result) *DBSpinResult {
	return &DBSpinResult{
		ID:        result.ID,
		UserID:    result.UserID,
		BetAmount: result.BetAmount,
		WinAmount: result.WinAmount,
		Reel1:     result.Reel1,
		Reel2:     result.Reel2,
		Reel3:     result.Reel3,
		IsWin:     result.IsWin,
		CreatedAt: result.CreatedAt,
	}
}

func toDomainSpinResult(dbResult *DBSpinResult) *spin.Result {
	return &spin.Result{
		ID:        dbResult.ID,
		UserID:    dbResult.UserID,
		BetAmount: dbResult.BetAmount,
		WinAmount: dbResult.WinAmount,
		Reel1:     dbResult.Reel1,
		Reel2:     dbResult.Reel2,
		Reel3:     dbResult.Reel3,
		IsWin:     dbResult.IsWin,
		CreatedAt: dbResult.CreatedAt,
	}
}

