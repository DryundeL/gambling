package spin

// Repository определяет интерфейс для работы с результатами спинов
type Repository interface {
	Create(result *Result) error
	GetByUserID(userID uint, limit int) ([]*Result, error)
}

