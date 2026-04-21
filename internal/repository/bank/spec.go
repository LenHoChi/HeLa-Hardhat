package bank

import (
	"context"

	"hela-bank-sc/internal/domain"
)

type Repository interface {
	Create(ctx context.Context, address string, action string, amount string, txHash string, status string) (*domain.History, error)
	ListByAddress(ctx context.Context, address string) ([]*domain.History, error)
}
