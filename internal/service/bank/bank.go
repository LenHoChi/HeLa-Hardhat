package bank

import (
	bank "hela-bank-sc/internal/blockchain"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

func (s *impl) GetBalance(addr common.Address) (*big.Int, error) {
	return bank.GetBalance(addr)
}

func (s *impl) Deposit(amount float64) (common.Hash, error) {
	return bank.Deposit(amount)
}

func (s *impl) Withdraw(amount float64) (common.Hash, error) {
	return bank.Withdraw(amount)
}

func (s *impl) EmergencyWithdraw() (common.Hash, error) {
	return bank.EmergencyWithdraw()
}

func (s *impl) GetContractBalance() (*big.Int, error) {
	return bank.GetContractBalance()
}
