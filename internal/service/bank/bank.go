package bank

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"hela-bank-sc/internal/blockchain"
	"hela-bank-sc/internal/models"
)

func (s impl) GetBalance(addr common.Address) (*big.Int, error) {
	return s.chain.GetBalance(addr)
}

func (s impl) Deposit(ctx context.Context, amount float64) (common.Hash, error) {
	txHash, amountWei, err := s.chain.Deposit(amount)
	if err != nil {
		return common.Hash{}, err
	}

	err = s.txRepo.Create(ctx,
		s.chain.FromAddress(),
		"deposit",
		amountWei.String(),
		txHash.Hex(),
		"submitted",
	)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (s impl) Withdraw(ctx context.Context, amount float64) (common.Hash, error) {
	txHash, amountWei, err := s.chain.Withdraw(amount)
	if err != nil {
		return common.Hash{}, err
	}

	err = s.txRepo.Create(ctx,
		blockchain.FromAddr.Hex(),
		"withdraw",
		amountWei.String(),
		txHash.Hex(),
		"submitted",
	)
	if err != nil {
		return common.Hash{}, err
	}

	return txHash, nil
}

func (s impl) EmergencyWithdraw(ctx context.Context) (common.Hash, error) {
	contractBalance, err := s.chain.GetContractBalance()
	if err != nil {
		return common.Hash{}, err
	}

	txHash, err := s.chain.EmergencyWithdraw()
	if err != nil {
		return common.Hash{}, err
	}

	err = s.txRepo.Create(ctx,
		blockchain.FromAddr.Hex(),
		"emergency_withdraw",
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
	return s.chain.GetContractBalance()
}

func (s impl) GetHistory(ctx context.Context, address string) ([]*models.TransactionHistory, error) {
	return s.txRepo.ListByAddress(ctx, address)
}
