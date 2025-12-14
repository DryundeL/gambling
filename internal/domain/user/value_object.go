package user

import "errors"

var (
	ErrInvalidCredentials = errors.New("неверные учетные данные")
)

// Credentials представляет Value Object для учетных данных
// Value Object - это объект без идентичности, определяемый только своими атрибутами
type Credentials struct {
	Username string
	Password string
}

// NewCredentials создает новый Value Object для учетных данных
func NewCredentials(username, password string) (*Credentials, error) {
	if username == "" || password == "" {
		return nil, ErrInvalidCredentials
	}
	return &Credentials{
		Username: username,
		Password: password,
	}, nil
}

