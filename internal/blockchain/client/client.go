package client

import (
	"fmt"
	"hela-bank-sc/internal/config"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	Client       *ethclient.Client
	ParsedABI    abi.ABI
	ContractAddr common.Address
)

func Init(cfg *config.Config) {
	var err error
	Client, err = ethclient.Dial(cfg.HELA_TESTNET_RPC)
	if err != nil {
		log.Fatal("Cannot connect to RPC:", err)
	} else {
		log.Println("Connected to RPC")
	}

	abiFile, err := os.ReadFile("abi/bank.json")
	if err != nil {
		log.Fatal("Cannot read ABI:", err)
	}
	ParsedABI, err = abi.JSON(strings.NewReader(string(abiFile)))
	if err != nil {
		log.Fatal("Cannot parse ABI:", err)
	}

	ContractAddr = common.HexToAddress(cfg.ContractAddress)
}

func GetClient() (*ethclient.Client, error) {
	if Client == nil {
		return nil, fmt.Errorf("RPC client is not initialized")
	}
	return Client, nil
}
