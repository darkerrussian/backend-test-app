package logic

import (
	"backend-test-app/internal/entity"
	"context"
	"errors"
	"testing"

	"github.com/shopspring/decimal"
)

// mockUserRepo — ручной мок entity.UserRepository.
type mockUserRepo struct {
	history entity.BalanceHistory
	err     error

	// calledWithUserID/calledWithAmount — чтобы проверить, что репозиторий
	// вызывается с ожидаемыми аргументами (полезно для теста на amount <= 0,
	// где репозиторий вообще не должен вызываться).
	called       bool
	calledUserID int64
	calledAmount decimal.Decimal
}

func (m *mockUserRepo) Withdraw(_ context.Context, userID int64, amount decimal.Decimal) (entity.BalanceHistory, error) {
	m.called = true
	m.calledUserID = userID
	m.calledAmount = amount

	if m.err != nil {
		return entity.BalanceHistory{}, m.err
	}
	return m.history, nil
}

func TestService_Withdraw_Success(t *testing.T) {
	expected := entity.BalanceHistory{
		UserId:     1,
		OldBalance: decimal.NewFromFloat(100),
		NewBalance: decimal.NewFromFloat(70),
		Amount:     decimal.NewFromFloat(30),
	}
	repo := &mockUserRepo{history: expected}
	service := NewUserLogic(repo)

	got, err := service.Withdraw(context.Background(), 1, decimal.NewFromFloat(30))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !got.NewBalance.Equal(expected.NewBalance) {
		t.Errorf("expected new balance %s, got %s", expected.NewBalance, got.NewBalance)
	}
	if !repo.called {
		t.Error("expected repo.Withdraw to be called")
	}
	if repo.calledUserID != 1 {
		t.Errorf("expected userID 1, got %d", repo.calledUserID)
	}
}

func TestService_Withdraw_InvalidAmount(t *testing.T) {
	repo := &mockUserRepo{}
	service := NewUserLogic(repo)

	_, err := service.Withdraw(context.Background(), 1, decimal.NewFromFloat(-10))
	if !errors.Is(err, entity.ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
	if repo.called {
		t.Error("repo.Withdraw should not be called for invalid amount")
	}
}

func TestService_Withdraw_ZeroAmount(t *testing.T) {
	repo := &mockUserRepo{}
	service := NewUserLogic(repo)

	_, err := service.Withdraw(context.Background(), 1, decimal.Zero)
	if !errors.Is(err, entity.ErrInvalidAmount) {
		t.Errorf("expected ErrInvalidAmount, got %v", err)
	}
}

func TestService_Withdraw_InsufficientBalance(t *testing.T) {
	repo := &mockUserRepo{err: entity.ErrInsufficientBalance}
	service := NewUserLogic(repo)

	_, err := service.Withdraw(context.Background(), 1, decimal.NewFromFloat(1000))
	if !errors.Is(err, entity.ErrInsufficientBalance) {
		t.Errorf("expected ErrInsufficientBalance, got %v", err)
	}
}

func TestService_Withdraw_UserNotFound(t *testing.T) {
	repo := &mockUserRepo{err: entity.ErrUserNotFound}
	service := NewUserLogic(repo)

	_, err := service.Withdraw(context.Background(), 999, decimal.NewFromFloat(10))
	if !errors.Is(err, entity.ErrUserNotFound) {
		t.Errorf("expected ErrUserNotFound, got %v", err)
	}
}
