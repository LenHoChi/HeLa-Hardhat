package bank

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"hela-bank-sc/internal/domain"
)

type Service interface {
	GetBalance(ctx context.Context, addr common.Address) (*big.Int, error)
	Deposit(ctx context.Context, amount float64) (common.Hash, error)
	Withdraw(ctx context.Context, amount float64) (common.Hash, error)
	EmergencyWithdraw(ctx context.Context) (common.Hash, error)
	GetContractBalance(ctx context.Context) (*big.Int, error)
	GetHistory(ctx context.Context, address string) ([]*domain.History, error)
}
