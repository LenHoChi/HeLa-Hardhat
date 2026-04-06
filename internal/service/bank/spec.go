package bank

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Service interface {
	GetBalance(addr common.Address) (*big.Int, error)
	Deposit(amount float64) (common.Hash, error)
	Withdraw(amount float64) (common.Hash, error)
}
