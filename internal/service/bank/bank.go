package bank

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"

	"hela-bank-sc/internal/domain"
)

func (s impl) GetBalance(ctx context.Context, addr common.Address) (*big.Int, error) {
	balance, err := s.chain.GetBalance(ctx, addr)
	if err != nil {
		return nil, fmt.Errorf("get balance from blockchain: %w", err)
	}

	return balance, nil
}

func (s impl) Deposit(ctx context.Context, amount float64) (common.Hash, error) {
	txHash, amountWei, err := s.chain.Deposit(ctx, amount)
	if err != nil {
		return common.Hash{}, fmt.Errorf("submit deposit transaction: %w", err)
	}

	err = s.txRepo.Create(ctx,
		s.chain.FromAddress(),
		"deposit",
		amountWei.String(),
		txHash.Hex(),
		"submitted",
	)
	if err != nil {
		return common.Hash{}, fmt.Errorf("create deposit history: %w", err)
	}

	return txHash, nil
}

func (s impl) Withdraw(ctx context.Context, amount float64) (common.Hash, error) {
	txHash, amountWei, err := s.chain.Withdraw(ctx, amount)
	if err != nil {
		return common.Hash{}, fmt.Errorf("submit withdraw transaction: %w", err)
	}

	err = s.txRepo.Create(ctx,
		s.chain.FromAddress(),
		"withdraw",
		amountWei.String(),
		txHash.Hex(),
		"submitted",
	)
	if err != nil {
		return common.Hash{}, fmt.Errorf("create withdraw history: %w", err)
	}

	return txHash, nil
}

func (s impl) EmergencyWithdraw(ctx context.Context) (common.Hash, error) {
	contractBalance, err := s.chain.GetContractBalance(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("get contract balance before emergency withdraw: %w", err)
	}

	txHash, err := s.chain.EmergencyWithdraw(ctx)
	if err != nil {
		return common.Hash{}, fmt.Errorf("submit emergency withdraw transaction: %w", err)
	}

	err = s.txRepo.Create(ctx,
		s.chain.FromAddress(),
		"emergency_withdraw",
		contractBalance.String(),
		txHash.Hex(),
		"submitted",
	)
	if err != nil {
		return common.Hash{}, fmt.Errorf("create emergency withdraw history: %w", err)
	}

	return txHash, nil
}

func (s impl) GetContractBalance(ctx context.Context) (*big.Int, error) {
	balance, err := s.chain.GetContractBalance(ctx)
	if err != nil {
		return nil, fmt.Errorf("get contract balance from blockchain: %w", err)
	}

	return balance, nil
}

func (s impl) GetHistory(ctx context.Context, address string) ([]*domain.History, error) {
	histories, err := s.txRepo.ListByAddress(ctx, address)
	if err != nil {
		return nil, fmt.Errorf("list history by address: %w", err)
	}

	return histories, nil
}
