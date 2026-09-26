package logic

import (
	"backend-test-app/internal/entity"
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

type userLogic struct {
	userRepository entity.UserRepository
}

func NewUserLogic(userRepository entity.UserRepository) entity.UserLogic {
	return &userLogic{
		userRepository: userRepository,
	}
}

func (l *userLogic) Withdraw(ctx context.Context, userId int64, amount decimal.Decimal) (entity.BalanceHistory, error) {
	if amount.LessThanOrEqual(decimal.Zero) {
		return entity.BalanceHistory{}, entity.ErrInvalidAmount
	}
	history, err := l.userRepository.Withdraw(ctx, userId, amount)
	if err != nil {
		return history, fmt.Errorf("withdraw: %w", err)
	}
	return history, nil
}
