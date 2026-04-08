package bank

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Service interface {
	GetBalance(addr common.Address) (*big.Int, error)
	Deposit(ctx context.Context, amount float64) (common.Hash, error)
	Withdraw(ctx context.Context, amount float64) (common.Hash, error)
	EmergencyWithdraw(ctx context.Context) (common.Hash, error)
	GetContractBalance() (*big.Int, error)
}
