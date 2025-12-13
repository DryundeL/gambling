package pgsql

import (
	"gambling/internal/model"
)

// RunMigrations выполняет автоматические миграции для всех моделей
func (s *Storage) RunMigrations() error {
	return s.DB.AutoMigrate(
		&model.User{},
		&model.Transaction{},
		&model.SpinResult{},
	)
}
