package pgsql

import (
	"fmt"
	"gambling/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Storage содержит экземпляр *gorm.DB для взаимодействия с базой данных.
type Storage struct {
	DB *gorm.DB
}

// New инициализирует новое подключение к базе данных с использованием GORM.
// Принимает структуру Config, содержащую параметры подключения.
func New(cfg *config.Config) *Storage {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=UTC",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
		cfg.DBSSLMode,
	)

	gormConfig := &gorm.Config{}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		panic("failed to connect to database using GORM:" + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get generic database object: %w" + err.Error())
	}

	if err := sqlDB.Ping(); err != nil {
		panic("failed to ping database: %w" + err.Error())
	}

	return &Storage{
		DB: db,
	}
}

// Close закрывает соединение с базой данных.
// Этот метод должен вызываться при завершении работы приложения.
func (s *Storage) Close() error {
	sqlDB, err := s.DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get generic database object: %w", err)
	}
	return sqlDB.Close()
}
