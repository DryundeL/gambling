package user

import "errors"

var (
	ErrInvalidAmount     = errors.New("неверная сумма")
	ErrInsufficientFunds = errors.New("недостаточно средств")
	ErrUserNotFound      = errors.New("пользователь не найден")
	ErrUserAlreadyExists = errors.New("пользователь уже существует")
)

