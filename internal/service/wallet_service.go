package service

import (
	"TestProject_itk/internal/domain"
	"context"

	"github.com/google/uuid"
)

type WalletRepository interface {
	Get(ctx context.Context, id uuid.UUID) (domain.Wallet, error)
	Deposit(ctx context.Context, id uuid.UUID, amount int64) (int64, error)
	Withdraw(ctx context.Context, id uuid.UUID, amount int64) (int64, error)
}

type WalletService struct {
	repo WalletRepository
}

func NewWalletService(repo WalletRepository) *WalletService {
	return &WalletService{repo: repo}
}

func (s *WalletService) GetBalance(ctx context.Context, id uuid.UUID) (int64, error) {
	w, err := s.repo.Get(ctx, id)
	if err != nil {
		return 0, err
	}
	return w.Balance, nil
}

func (s *WalletService) ApplyOperation(ctx context.Context, id uuid.UUID, op domain.OperationType, amount int64) (int64, error) {
	if amount <= 0 {
		return 0, domain.ErrInsufficientFund
	}

	switch op {
	case domain.OperationDeposit:
		return s.repo.Deposit(ctx, id, amount)
	case domain.OperationWithdraw:
		return s.repo.Withdraw(ctx, id, amount)
	default:
		return 0, domain.ErrWalletNotFound
	}
}
