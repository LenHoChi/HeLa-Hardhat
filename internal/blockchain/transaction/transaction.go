package transaction

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	clientpkg "hela-bank-sc/internal/blockchain/client"
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
	PrivateKey, err = crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		log.Fatal("Cannot load private key:", err)
	}
	publicKey := PrivateKey.Public().(*ecdsa.PublicKey)
	FromAddr = crypto.PubkeyToAddress(*publicKey)
}

func GetAuth() (*bind.TransactOpts, error) {
	nonce, err := clientpkg.Client.PendingNonceAt(context.Background(), FromAddr)
	if err != nil {
		log.Fatal("Cannot get nonce:", err)
	}
	gasPrice, err := clientpkg.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return nil, fmt.Errorf("get gas price: %w", err)
	}
	auth, err := bind.NewKeyedTransactorWithChainID(PrivateKey, big.NewInt(666888))
	if err != nil {
		return nil, fmt.Errorf("create transactor: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasPrice = gasPrice
	auth.GasLimit = uint64(100000)
	return auth, nil
}

func WaitForTx(txHash common.Hash) {
	fmt.Println("⏳ Waiting for transaction confirm...")
	for {
		receipt, err := clientpkg.Client.TransactionReceipt(context.Background(), txHash)
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

func FromAddress() string {
	return FromAddr.Hex()
}
