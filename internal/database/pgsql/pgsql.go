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
		fmt.Println("═══════════════════════════════════════════════════════════")
		fmt.Println("❌ ОШИБКА ПОДКЛЮЧЕНИЯ К БАЗЕ ДАННЫХ")
		fmt.Println("═══════════════════════════════════════════════════════════")
		fmt.Printf("Не удалось подключиться к PostgreSQL: %v\n", err)
		fmt.Println("")
		fmt.Println("Проверьте:")
		fmt.Println("  1. PostgreSQL запущен и доступен")
		fmt.Println("  2. Параметры в файле .env корректны:")
		fmt.Printf("     - DB_HOST=%s\n", cfg.DBHost)
		fmt.Printf("     - DB_PORT=%d\n", cfg.DBPort)
		fmt.Printf("     - DB_USER=%s\n", cfg.DBUser)
		fmt.Printf("     - DB_NAME=%s\n", cfg.DBName)
		fmt.Println("     - DB_PASSWORD=*** (проверьте пароль)")
		fmt.Println("  3. База данных существует:")
		fmt.Printf("     CREATE DATABASE %s;\n", cfg.DBName)
		fmt.Println("  4. Пользователь имеет права на базу данных")
		fmt.Println("═══════════════════════════════════════════════════════════")
		panic("failed to connect to database using GORM:" + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("failed to get generic database object: %w" + err.Error())
	}

	if err := sqlDB.Ping(); err != nil {
		panic("failed to ping database: %w" + err.Error())
	}

	storage := &Storage{
		DB: db,
	}

	// Выполняем миграции при инициализации
	if err := storage.RunMigrations(); err != nil {
		panic("failed to run migrations: " + err.Error())
	}

	return storage
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
