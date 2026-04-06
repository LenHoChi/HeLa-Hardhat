package bank

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"log"
	"math/big"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	PrivateKey *ecdsa.PrivateKey
	FromAddr   common.Address
)

func InitWallet() {
	var err error
	// PRIVATE_KEY=f3484531ac9e7c578ef47a356147af525f1234ad6b957a5dccce3b6a6a77a0f9 (get in ur wallet)
	PrivateKey, err = crypto.HexToECDSA(os.Getenv("PRIVATE_KEY")) // load key to sign txn
	if err != nil {
		log.Fatal("Cannot load private key:", err)
	}
	publicKey := PrivateKey.Public().(*ecdsa.PublicKey) // get public key from private key
	FromAddr = crypto.PubkeyToAddress(*publicKey)       // derive wallet address from public key
	// private key -> public key -> wallet address (0x32A413fc36E202849B4eDffdB111802804fC7AEe)
}

func GetAuth() *bind.TransactOpts { // return TransactOpts (include all info to send txn to blockchain)
	nonce, err := Client.PendingNonceAt(context.Background(), FromAddr) // get nonce
	if err != nil {
		log.Fatal("Cannot get nonce:", err)
	}
	gasPrice, err := Client.SuggestGasPrice(context.Background()) // get gas
	if err != nil {
		log.Fatal("Cannot get gas price:", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(PrivateKey, big.NewInt(666888)) // create object to sign txn
	if err != nil {
		log.Fatal("Cannot create transactor:", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(100000)
	return auth
}

func WaitForTx(txHash common.Hash) {
	fmt.Println("⏳ Waiting for transaction confirm...")
	for {
		receipt, err := Client.TransactionReceipt(context.Background(), txHash)
		if err == nil && receipt != nil {
			if receipt.Status == 1 {
				fmt.Println("✅ Transaction confirmed!")
			} else {
				fmt.Println("❌ Transaction failed!")
			}
			return
		}
		time.Sleep(2 * time.Second)
	}
}
