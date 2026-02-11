package repository

import (
	"context"

	"TestProject_itk/internal/domain"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type WalletRepo struct {
	pool *pgxpool.Pool
}

func NewWalletRepo(pool *pgxpool.Pool) *WalletRepo {
	return &WalletRepo{pool: pool}
}

func (r *WalletRepo) Get(ctx context.Context, id uuid.UUID) (domain.Wallet, error) {
	var w domain.Wallet
	w.ID = id

	err := r.pool.QueryRow(ctx, `SELECT balance FROM wallets WHERE id=$1`, id).Scan(&w.Balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return domain.Wallet{}, domain.ErrWalletNotFound
		}
		return domain.Wallet{}, err
	}

	return w, nil
}

func (r *WalletRepo) Deposit(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	var balance int64

	_, err := r.pool.Exec(ctx, `
		INSERT INTO wallets (id, balance)
		VALUES ($1, 0)
		ON CONFLICT (id) DO NOTHING
	`, id)
	if err != nil {
		return 0, err
	}

	err = r.pool.QueryRow(ctx, `
		UPDATE wallets
		SET balance = balance + $2
		WHERE id = $1
		RETURNING balance
	`, id, amount).Scan(&balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			return 0, domain.ErrWalletNotFound
		}
		return 0, err
	}

	return balance, nil
}

func (r *WalletRepo) Withdraw(ctx context.Context, id uuid.UUID, amount int64) (int64, error) {
	var balance int64

	err := r.pool.QueryRow(ctx, `
		UPDATE wallets
		SET balance = balance - $2
		WHERE id = $1 AND balance >= $2
		RETURNING balance
	`, id, amount).Scan(&balance)
	if err != nil {
		if err == pgx.ErrNoRows {
			var exists bool
			err2 := r.pool.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM wallets WHERE id=$1)`, id).Scan(&exists)
			if err2 != nil {
				return 0, err2
			}
			if !exists {
				return 0, domain.ErrWalletNotFound
			}
			return 0, domain.ErrInsufficientFund
		}
		return 0, err
	}

	return balance, nil
}
