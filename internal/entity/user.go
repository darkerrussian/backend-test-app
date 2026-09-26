package entity

import (
	"context"
	"errors"
	"time"

	"github.com/shopspring/decimal"
)

var (
	// ErrUserNotFound возвращается, если пользователь с указанным ID не найден.
	ErrUserNotFound = errors.New("user not found")
	// ErrInsufficientBalance возвращается, когда баланса не хватает на списание.
	ErrInsufficientBalance = errors.New("insufficient balance")
	// ErrInvalidAmount возвращается для суммы списания <= 0.
	ErrInvalidAmount = errors.New("amount must be positive")
)

// User - пользователь с балансом
type User struct {
	Id int64
	// используем decimal вместо float64 чтобы не было неточностей при округлении
	Balance decimal.Decimal
}

type BalanceHistory struct {
	Id         int64
	UserId     int64
	OldBalance decimal.Decimal
	NewBalance decimal.Decimal
	Amount     decimal.Decimal
	CreatedAt  time.Time
}

// UserRepository - то, что нужно от хранилища для списания баланса.
// Withdraw обязана быть атомарной: обновление баланса и запись истории
// должны происходить в одной транзакции.
type UserRepository interface {
	Withdraw(ctx context.Context, userId int64, amount decimal.Decimal) (BalanceHistory, error)
}

type UserLogic interface {
	Withdraw(ctx context.Context, userId int64, amount decimal.Decimal) (BalanceHistory, error)
}
