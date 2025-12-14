package pgsql

import (
	"gambling/internal/infrastructure/repository"
)

// RunMigrations выполняет автоматические миграции для всех моделей
func (s *Storage) RunMigrations() error {
	return s.DB.AutoMigrate(
		&repository.DBUser{},
		&repository.DBTransaction{},
		&repository.DBSpinResult{},
	)
}

