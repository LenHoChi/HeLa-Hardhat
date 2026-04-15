package blockchain

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Gateway interface {
	GetBalance(ctx context.Context, addr common.Address) (*big.Int, error)
	Deposit(ctx context.Context, amount float64) (common.Hash, *big.Int, error)
	Withdraw(ctx context.Context, amount float64) (common.Hash, *big.Int, error)
	EmergencyWithdraw(ctx context.Context) (common.Hash, error)
	GetContractBalance(ctx context.Context) (*big.Int, error)
	FromAddress() string
}
