package blockchain

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Gateway interface {
	GetBalance(addr common.Address) (*big.Int, error)
	Deposit(amount float64) (common.Hash, *big.Int, error)
	Withdraw(amount float64) (common.Hash, *big.Int, error)
	EmergencyWithdraw() (common.Hash, error)
	GetContractBalance() (*big.Int, error)
	FromAddress() string
}
