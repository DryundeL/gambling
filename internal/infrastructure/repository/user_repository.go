package repository

import (
	"errors"
	"gambling/internal/domain/user"
	"time"

	"gorm.io/gorm"
)

// UserRepository реализует интерфейс user.Repository
// Это адаптер (adapter) в архитектуре Ports & Adapters
// Он адаптирует доменный интерфейс к конкретной реализации БД (GORM)
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository создает новый репозиторий пользователей
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// Create создает нового пользователя
func (r *UserRepository) Create(u *user.User) error {
	// Преобразуем доменную сущность в модель для БД
	dbUser := toDBModel(u)
	if err := r.db.Create(dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return user.ErrUserAlreadyExists
		}
		return err
	}
	// Обновляем ID в доменной сущности
	u.ID = dbUser.ID
	u.CreatedAt = dbUser.CreatedAt
	u.UpdatedAt = dbUser.UpdatedAt
	return nil
}

// GetByID возвращает пользователя по ID
func (r *UserRepository) GetByID(id uint) (*user.User, error) {
	var dbUser DBUser
	if err := r.db.First(&dbUser, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return toDomainModel(&dbUser), nil
}

// GetByUsername возвращает пользователя по имени пользователя
func (r *UserRepository) GetByUsername(username string) (*user.User, error) {
	var dbUser DBUser
	if err := r.db.Where("username = ?", username).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return toDomainModel(&dbUser), nil
}

// GetByEmail возвращает пользователя по email
func (r *UserRepository) GetByEmail(email string) (*user.User, error) {
	var dbUser DBUser
	if err := r.db.Where("email = ?", email).First(&dbUser).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, user.ErrUserNotFound
		}
		return nil, err
	}
	return toDomainModel(&dbUser), nil
}

// UpdateBalance обновляет баланс пользователя
func (r *UserRepository) UpdateBalance(userID uint, newBalance float64) error {
	return r.db.Model(&DBUser{}).Where("id = ?", userID).Update("balance", newBalance).Error
}

// Update обновляет данные пользователя
func (r *UserRepository) Update(u *user.User) error {
	dbUser := toDBModel(u)
	return r.db.Save(dbUser).Error
}

// DBUser представляет модель БД для пользователя
// Это техническая деталь инфраструктуры, отделенная от домена
type DBUser struct {
	ID           uint           `gorm:"primaryKey"`
	Username     string         `gorm:"uniqueIndex;not null;size:50"`
	Email        string         `gorm:"uniqueIndex;not null;size:100"`
	PasswordHash string         `gorm:"not null;size:255"`
	Balance      float64        `gorm:"not null;default:0;type:decimal(15,2)"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}

func (DBUser) TableName() string {
	return "users"
}

// toDBModel преобразует доменную сущность в модель БД
func toDBModel(u *user.User) *DBUser {
	return &DBUser{
		ID:           u.ID,
		Username:     u.Username,
		Email:        u.Email,
		PasswordHash: u.PasswordHash,
		Balance:      u.Balance,
	}
}

// toDomainModel преобразует модель БД в доменную сущность
func toDomainModel(dbUser *DBUser) *user.User {
	return &user.User{
		ID:           dbUser.ID,
		Username:     dbUser.Username,
		Email:        dbUser.Email,
		PasswordHash: dbUser.PasswordHash,
		Balance:      dbUser.Balance,
		CreatedAt:    dbUser.CreatedAt,
		UpdatedAt:    dbUser.UpdatedAt,
		DeletedAt:    dbUser.DeletedAt,
	}
}

