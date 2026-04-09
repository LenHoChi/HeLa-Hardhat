package blockchain

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

func (g gateway) GetBalance(userAddr common.Address) (*big.Int, error) {
	data, err := ParsedABI.Pack("getBalance", userAddr) // paste address to check balance
	if err != nil {
		return nil, err
	}
	result, err := Client.CallContract(context.Background(), ethereum.CallMsg{
		To:   &ContractAddr,
		Data: data,
	}, nil)
	if err != nil {
		return nil, err
	}
	balance := new(big.Int).SetBytes(result)
	return balance, nil
}

func (g gateway) Deposit(amountEther float64) (common.Hash, *big.Int, error) {
	auth := GetAuth()
	amount := new(big.Float).Mul(
		big.NewFloat(amountEther),
		big.NewFloat(1e18),
	)
	amountWei, _ := amount.Int(nil)
	auth.Value = amountWei

	data, err := ParsedABI.Pack("deposit")
	if err != nil {
		return common.Hash{}, nil, err
	}
	tx, err := bind.NewBoundContract(ContractAddr, ParsedABI, Client, Client, Client).
		RawTransact(auth, data)
	if err != nil {
		return common.Hash{}, nil, err
	}
	fmt.Printf("✅ Deposit %.2f ETH — tx: %s\n", amountEther, tx.Hash().Hex())
	return tx.Hash(), amountWei, nil
}

func (g gateway) Withdraw(amountEther float64) (common.Hash, *big.Int, error) {
	auth := GetAuth()
	amount := new(big.Float).Mul(
		big.NewFloat(amountEther),
		big.NewFloat(1e18),
	)
	amountWei, _ := amount.Int(nil)

	data, err := ParsedABI.Pack("withdraw", amountWei)
	if err != nil {
		return common.Hash{}, nil, err
	}
	tx, err := bind.NewBoundContract(ContractAddr, ParsedABI, Client, Client, Client).
		RawTransact(auth, data)
	if err != nil {
		return common.Hash{}, nil, err
	}
	fmt.Printf("✅ Withdraw %.2f ETH — tx: %s\n", amountEther, tx.Hash().Hex())
	return tx.Hash(), amountWei, nil
}

func (g gateway) EmergencyWithdraw() (common.Hash, error) {
	auth := GetAuth()

	data, err := ParsedABI.Pack("emergencyWithdraw")
	if err != nil {
		return common.Hash{}, err
	}
	tx, err := bind.NewBoundContract(ContractAddr, ParsedABI, Client, Client, Client).
		RawTransact(auth, data)
	if err != nil {
		return common.Hash{}, err
	}
	return tx.Hash(), nil
}

func (g gateway) GetContractBalance() (*big.Int, error) {
	balance, err := Client.BalanceAt(context.Background(), ContractAddr, nil)
	if err != nil {
		return nil, err
	}
	return balance, nil
}

func (g gateway) FromAddress() string {
	return FromAddr.Hex()
}

func PrintBalance(userAddr common.Address) {
	gw := gateway{}
	balance, err := gw.GetBalance(userAddr)
	if err != nil {
		log.Println("Cannot get balance:", err)
		return
	}
	fmt.Printf("✅ Balance of %s: %s wei\n", userAddr.Hex(), balance.String())
}
