package blockchain

import (
	"context"
	"fmt"
	clientpkg "hela-bank-sc/internal/blockchain/client"
	txpkg "hela-bank-sc/internal/blockchain/transaction"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (g gateway) GetBalance(ctx context.Context, userAddr common.Address) (*big.Int, error) {
	client, err := clientpkg.GetClient()
	if err != nil {
		return nil, fmt.Errorf("get RPC client: %w", err)
	}

	data, err := clientpkg.ParsedABI.Pack("getBalance", userAddr) // paste address to check balance
	if err != nil {
		return nil, err
	}
	result, err := client.CallContract(ctx, ethereum.CallMsg{
		To:   &clientpkg.ContractAddr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}
	balance := new(big.Int).SetBytes(result)
	return balance, nil
}

func (g gateway) Deposit(ctx context.Context, amountEther float64) (common.Hash, *big.Int, error) {
	client, err := clientpkg.GetClient()
	if err != nil {
		return common.Hash{}, nil, fmt.Errorf("get RPC client: %w", err)
	}

	auth, err := txpkg.GetAuth(ctx) // load private key
	if err != nil {
		return common.Hash{}, nil, err
	}

	amount := new(big.Float).Mul(
		big.NewFloat(amountEther),
		big.NewFloat(1e18),
	)
	amountWei, _ := amount.Int(nil)
	auth.Value = amountWei

	data, err := clientpkg.ParsedABI.Pack("deposit")
	if err != nil {
		return common.Hash{}, nil, err
	}
	tx, err := bind.NewBoundContract(clientpkg.ContractAddr, clientpkg.ParsedABI, client, client, client).
		RawTransact(auth, data) // ký và broadcast
	if err != nil {
		return common.Hash{}, nil, err
	}
	fmt.Printf("✅ Deposit %.2f ETH — tx: %s\n", amountEther, tx.Hash().Hex())
	return tx.Hash(), amountWei, nil
}

func (g gateway) Withdraw(ctx context.Context, amountEther float64) (common.Hash, *big.Int, error) {
	client, err := clientpkg.GetClient()
	if err != nil {
		return common.Hash{}, nil, fmt.Errorf("get RPC client: %w", err)
	}

	auth, err := txpkg.GetAuth(ctx)
	if err != nil {
		return common.Hash{}, nil, err
	}

	amount := new(big.Float).Mul(
		big.NewFloat(amountEther),
		big.NewFloat(1e18),
	)
	amountWei, _ := amount.Int(nil)

	data, err := clientpkg.ParsedABI.Pack("withdraw", amountWei)
	if err != nil {
		return common.Hash{}, nil, err
	}
	tx, err := bind.NewBoundContract(clientpkg.ContractAddr, clientpkg.ParsedABI, client, client, client).
		RawTransact(auth, data)
	if err != nil {
		return common.Hash{}, nil, err
	}
	fmt.Printf("✅ Withdraw %.2f ETH — tx: %s\n", amountEther, tx.Hash().Hex())
	return tx.Hash(), amountWei, nil
}

func (g gateway) EmergencyWithdraw(ctx context.Context) (common.Hash, error) {
	client, err := clientpkg.GetClient()
	if err != nil {
		return common.Hash{}, fmt.Errorf("get RPC client: %w", err)
	}

	auth, err := txpkg.GetAuth(ctx)
	if err != nil {
		return common.Hash{}, err
	}

	data, err := clientpkg.ParsedABI.Pack("emergencyWithdraw")
	if err != nil {
		return common.Hash{}, err
	}
	tx, err := bind.NewBoundContract(clientpkg.ContractAddr, clientpkg.ParsedABI, client, client, client).
		RawTransact(auth, data)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

func (g gateway) GetContractBalance(ctx context.Context) (*big.Int, error) {
	client, err := clientpkg.GetClient()
	if err != nil {
		return nil, fmt.Errorf("get RPC client: %w", err)
	}

	balance, err := client.BalanceAt(ctx, clientpkg.ContractAddr, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (g gateway) FromAddress() string {
	return txpkg.FromAddress()
}

func PrintBalance(userAddr common.Address) {
	gw := gateway{}
	balance, err := gw.GetBalance(context.Background(), userAddr)
	if err != nil {
		log.Println("Cannot get balance:", err)
		return
	}
	fmt.Printf("✅ Balance of %s: %s wei\n", userAddr.Hex(), balance.String())
}
