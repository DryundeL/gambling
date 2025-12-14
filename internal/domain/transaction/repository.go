package transaction

// Repository определяет интерфейс для работы с транзакциями
type Repository interface {
	Create(transaction *Transaction) error
	GetByUserID(userID uint, limit int) ([]*Transaction, error)
}

