package blockchain

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
	ContractAddress = "0x85933342B34ceB2ef5ECc63FEC7659c4a3495d6F"
	WssURL          = "wss://testnet-rpc.helachain.com"
)

const ()

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
