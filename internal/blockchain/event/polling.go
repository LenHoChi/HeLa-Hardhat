package event

import (
	"context"
	"fmt"
	clientpkg "hela-bank-sc/internal/blockchain/client"
	"log"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type DepositedEvent struct {
	User   common.Address
	Amount *big.Int
}

type WithdrawnEvent struct {
	User   common.Address
	Amount *big.Int
}

func ListenPolling(ctx context.Context) {
	fmt.Println("👂 Listening to events...")

	latestBlock, err := clientpkg.Client.BlockNumber(ctx)
	if err != nil {
		log.Fatal("Cannot get block number:", err)
		// log.Println("Cannot get block number: ", err)
	}

	fromBlock := latestBlock

	for {
		select {
		case <-ctx.Done():
			fmt.Println("🛑 Stop listening events")
			return
		default:
			currentBlock, err := clientpkg.Client.BlockNumber(ctx)
			if err != nil {
				log.Println("Cannot get block number:", err)
				time.Sleep(3 * time.Second)
				continue
			}

			if currentBlock > fromBlock {
				query := ethereum.FilterQuery{
					Addresses: []common.Address{clientpkg.ContractAddr},
					FromBlock: big.NewInt(int64(fromBlock)),
					ToBlock:   big.NewInt(int64(currentBlock)),
				}

				logs, err := clientpkg.Client.FilterLogs(ctx, query)
				if err != nil {
					log.Println("Cannot filter logs:", err)
					time.Sleep(3 * time.Second)
					continue
				}

				for _, vLog := range logs {
					handleEvent(vLog)
				}

				fromBlock = currentBlock + 1
			}

			time.Sleep(3 * time.Second)
		}
	}
}

func handleEvent(vLog types.Log) {
	depositedSig := clientpkg.ParsedABI.Events["Deposited"].ID
	withdrawnSig := clientpkg.ParsedABI.Events["Withdrawn"].ID
	emergencySig := clientpkg.ParsedABI.Events["EmergencyWithdrawn"].ID

	switch vLog.Topics[0] {
	case depositedSig:
		handleDeposited(vLog)
	case withdrawnSig:
		handleWithdrawn(vLog)
	case emergencySig:
		handleEmergency(vLog)
	}
}

func handleDeposited(vLog types.Log) {
	var event DepositedEvent
	err := clientpkg.ParsedABI.UnpackIntoInterface(&event, "Deposited", vLog.Data)
	if err != nil {
		log.Println("Cannot unpack Deposited event:", err)
		return
	}
	event.User = common.HexToAddress(vLog.Topics[1].Hex())
	fmt.Printf("📥 Deposited — user: %s, amount: %s wei\n",
		event.User.Hex(), event.Amount.String())
}

func handleWithdrawn(vLog types.Log) {
	var event WithdrawnEvent
	err := clientpkg.ParsedABI.UnpackIntoInterface(&event, "Withdrawn", vLog.Data)
	if err != nil {
		log.Println("Cannot unpack Withdrawn event:", err)
		return
	}
	event.User = common.HexToAddress(vLog.Topics[1].Hex())
	fmt.Printf("📤 Withdrawn — user: %s, amount: %s wei\n",
		event.User.Hex(), event.Amount.String())
}

func handleEmergency(vLog types.Log) {
	var event struct {
		Amount *big.Int
	}
	err := clientpkg.ParsedABI.UnpackIntoInterface(&event, "EmergencyWithdrawn", vLog.Data)
	if err != nil {
		log.Println("Cannot unpack EmergencyWithdrawn event:", err)
		return
	}
	owner := common.HexToAddress(vLog.Topics[1].Hex())
	fmt.Printf("🚨 EmergencyWithdrawn — owner: %s, amount: %s wei\n",
		owner.Hex(), event.Amount.String())
}
