package postgres

import (
	"backend-test-app/internal/entity"
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/shopspring/decimal"
)

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) entity.UserRepository {
	return &userRepository{
		pool: pool,
	}
}

func (u userRepository) Withdraw(ctx context.Context, userId int64, amount decimal.Decimal) (entity.BalanceHistory, error) {
	tx, err := u.pool.Begin(ctx)
	if err != nil {
		return entity.BalanceHistory{}, fmt.Errorf("could not start transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	var oldBalance decimal.Decimal
	err = tx.QueryRow(ctx,
		`SELECT balance FROM users WHERE id = $1 FOR UPDATE`,
		userId).Scan(&oldBalance)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.BalanceHistory{}, entity.ErrUserNotFound
	}
	if err != nil {
		return entity.BalanceHistory{}, fmt.Errorf("could not query old balance: %w", err)
	}
	if oldBalance.LessThan(amount) {
		return entity.BalanceHistory{}, entity.ErrInsufficientBalance
	}

	newBalance := oldBalance.Sub(amount)

	_, err = tx.Exec(ctx,
		`UPDATE users SET balance = $1 WHERE id = $2`,
		newBalance, userId)
	if err != nil {
		return entity.BalanceHistory{}, fmt.Errorf("could not update balance: %w", err)
	}
	var history entity.BalanceHistory
	err = tx.QueryRow(ctx,
		`INSERT INTO balance_history (user_id, old_balance, new_balance, amount)
		 VALUES ($1, $2, $3, $4)
		 RETURNING id, user_id, old_balance, new_balance, amount, created_at`,
		userId, oldBalance, newBalance, amount,
	).Scan(&history.Id, &history.UserId, &history.OldBalance, &history.NewBalance, &history.Amount, &history.CreatedAt)
	if err != nil {
		return entity.BalanceHistory{}, fmt.Errorf("could not insert new balance: %w", err)
	}
	if err := tx.Commit(ctx); err != nil {
		return entity.BalanceHistory{}, fmt.Errorf("could not commit transaction: %w", err)
	}
	return history, nil
}
