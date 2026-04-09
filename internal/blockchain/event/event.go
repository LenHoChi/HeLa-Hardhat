package event

import (
	"context"
	"encoding/json"
	"fmt"
	clientpkg "hela-bank-sc/internal/blockchain/client"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const explorerAPI = "https://testnet-blockexplorer.helachain.com/api/v2/addresses/"

type ExplorerLog struct {
	BlockNumber int    `json:"block_number"`
	TxHash      string `json:"tx_hash"`
	Decoded     struct {
		MethodCall string `json:"method_call"`
		Parameters []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"parameters"`
	} `json:"decoded"`
}

type ExplorerResponse struct {
	Items []ExplorerLog `json:"items"`
}

func ListenExplorer(ctx context.Context) {
	fmt.Println("👂 Listening to events via Explorer API...")

	latestBlock, err := clientpkg.Client.BlockNumber(ctx)
	if err != nil {
		log.Println("Cannot get latest block:", err)
		latestBlock = 0
	}
	lastBlock := int(latestBlock)
	fmt.Printf("📌 Starting from block: %d\n", lastBlock)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("🛑 Stop listening events")
			return
		default:
			logs, err := fetchLogs()
			if err != nil {
				log.Println("Cannot fetch logs:", err)
				time.Sleep(3 * time.Second)
				continue
			}

			for _, item := range logs {
				if item.BlockNumber <= lastBlock {
					continue
				}
				handleExplorerEvent(item)
				lastBlock = item.BlockNumber
			}

			time.Sleep(3 * time.Second)
		}
	}
}

func fetchLogs() ([]ExplorerLog, error) {
	url := explorerAPI + clientpkg.ContractAddr.Hex() + "/logs"
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var result ExplorerResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	return result.Items, nil
}

func handleExplorerEvent(item ExplorerLog) {
	method := item.Decoded.MethodCall
	amount := getParam(item.Decoded.Parameters, "amount")

	switch {
	case strings.Contains(method, "EmergencyWithdrawn"):
		owner := getParam(item.Decoded.Parameters, "owner")
		fmt.Printf("🚨 EmergencyWithdrawn — owner: %s, amount: %s wei | tx: %s\n",
			owner, amount, item.TxHash)
	case strings.Contains(method, "Deposited"):
		user := getParam(item.Decoded.Parameters, "user")
		fmt.Printf("📥 Deposited — user: %s, amount: %s wei | tx: %s\n",
			user, amount, item.TxHash)
	case strings.Contains(method, "Withdrawn"):
		user := getParam(item.Decoded.Parameters, "user")
		fmt.Printf("📤 Withdrawn — user: %s, amount: %s wei | tx: %s\n",
			user, amount, item.TxHash)
	}
}

func getParam(params []struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}, name string) string {
	for _, p := range params {
		if p.Name == name {
			return p.Value
		}
	}
	return ""
}
