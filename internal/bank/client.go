package bank

import (
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/joho/godotenv"
)

const (
	RpcURL          = "https://666888.rpc.thirdweb.com"
	ContractAddress = "0x8a191B6D92FDB9A75EDb3cBeB6e78105d44B2822"
)

var (
	Client       *ethclient.Client
	ParsedABI    abi.ABI
	ContractAddr common.Address
)

func Init() {
	// Load .env
	godotenv.Load()

	// Kết nối RPC
	var err error
	Client, err = ethclient.Dial(RpcURL) // connect to HeLa testnet by rpcURL
	if err != nil {
		log.Fatal("Cannot connect to RPC:", err)
	}

	// Load ABI
	abiFile, err := os.ReadFile("abi/bank.json") // load ABI (ABI = all scenario deployed contract)
	if err != nil {
		log.Fatal("Cannot read ABI:", err)
	}
	ParsedABI, err = abi.JSON(strings.NewReader(string(abiFile))) // parse file to json
	if err != nil {
		log.Fatal("Cannot parse ABI:", err)
	}

	// Set contract address
	ContractAddr = common.HexToAddress(ContractAddress)
}
