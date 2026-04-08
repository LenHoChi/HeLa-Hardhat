package bank

import (
	"context"
	"hela-bank-sc/internal/blockchain"
	bank "hela-bank-sc/internal/blockchain"
	"math/big"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
)

func (s impl) GetBalance(addr common.Address) (*big.Int, error) {
	return bank.GetBalance(addr)
}

func (s impl) Deposit(ctx context.Context, amount float64) (common.Hash, error) {
	txHash, err := bank.Deposit(amount)
	if err != nil {
		return common.Hash{}, err
	}

	err = s.txRepo.Create(ctx,
		blockchain.FromAddr.Hex(),
		"deposit",
		strconv.FormatFloat(amount, 'f', -1, 64),
		txHash.Hex(),
		"submitted",
	)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (s impl) Withdraw(ctx context.Context, amount float64) (common.Hash, error) {
	txHash, err := bank.Withdraw(amount)
	if err != nil {
		return common.Hash{}, err
	}

	err = s.txRepo.Create(ctx,
		blockchain.FromAddr.Hex(),
		"withdraw",
		strconv.FormatFloat(amount, 'f', -1, 64),
		txHash.Hex(),
		"submitted",
	)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (s impl) EmergencyWithdraw(ctx context.Context) (common.Hash, error) {
	contractBalance, err := bank.GetContractBalance()
	if err != nil {
		return common.Hash{}, err
	}

	txHash, err := bank.EmergencyWithdraw()
	if err != nil {
		return common.Hash{}, err
	}

	err = s.txRepo.Create(ctx,
		blockchain.FromAddr.Hex(),
		"emergency withdraw",
		contractBalance.String(),
		txHash.Hex(),
		"submitted",
	)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (s impl) GetContractBalance() (*big.Int, error) {
	return bank.GetContractBalance()
}
