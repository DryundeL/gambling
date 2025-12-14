package user

// Repository определяет интерфейс для работы с пользователями
// Это порт (port) в архитектуре Ports & Adapters (Hexagonal Architecture)
// Реализация находится в infrastructure слое
type Repository interface {
	Create(user *User) error
	GetByID(id uint) (*User, error)
	GetByUsername(username string) (*User, error)
	GetByEmail(email string) (*User, error)
	UpdateBalance(userID uint, newBalance float64) error
	Update(user *User) error
}

