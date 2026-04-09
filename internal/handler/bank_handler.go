package handler

import (
	"encoding/json"
	bank "hela-bank-sc/internal/blockchain"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
)

type DepositRequest struct {
	Amount float64 `json:"amount"`
}

type WithdrawRequest struct {
	Amount float64 `json:"amount"`
}

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type HistoryItem struct {
	Address   string    `json:"address"`
	Action    string    `json:"action"`
	Amount    string    `json:"amount"`
	TxHash    string    `json:"tx_hash"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (h BankHandler) CheckHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, Response{
			Success: true,
			Message: "Check health ok",
		})
	}
}

func (h BankHandler) GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := chi.URLParam(r, "address")

		userAddr := common.HexToAddress(address)
		balance, err := h.bankSvc.GetBalance(userAddr)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, Response{
				Success: false,
				Message: "Cannot get balance: " + err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, Response{
			Success: true,
			Message: "Get balance successfully",
			Data: map[string]string{
				"address": address,
				"balance": balance.String(),
			},
		})
	}
}

func (h BankHandler) Deposit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DepositRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, Response{
				Success: false,
				Message: "Invalid request body",
			})
			return
		}

		txHash, err := h.bankSvc.Deposit(r.Context(), req.Amount)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, Response{
				Success: false,
				Message: "Cannot deposit: " + err.Error(),
			})
			return
		}

		go bank.WaitForTx(txHash)

		writeJSON(w, http.StatusOK, Response{
			Success: true,
			Message: "Deposit submitted",
			Data: map[string]string{
				"tx_hash":  txHash.Hex(),
				"explorer": "https://testnet-blockexplorer.helachain.com/tx/" + txHash.Hex(),
			},
		})
	}
}

func (h BankHandler) Withdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req WithdrawRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, Response{
				Success: false,
				Message: "Invalid request body",
			})
			return
		}

		txHash, err := h.bankSvc.Withdraw(r.Context(), req.Amount)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, Response{
				Success: false,
				Message: "Cannot withdraw: " + err.Error(),
			})
			return
		}

		go bank.WaitForTx(txHash)

		writeJSON(w, http.StatusOK, Response{
			Success: true,
			Message: "Withdraw submitted",
			Data: map[string]string{
				"tx_hash":  txHash.Hex(),
				"explorer": "https://testnet-blockexplorer.helachain.com/tx/" + txHash.Hex(),
			},
		})
	}
}

func (h BankHandler) EmergencyWithdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txHash, err := h.bankSvc.EmergencyWithdraw(r.Context())
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, Response{
				Success: false,
				Message: "Cannot emergency withdraw: " + err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, Response{
			Success: true,
			Message: "Emergency withdraw submitted",
			Data: map[string]string{
				"tx_hash":  txHash.Hex(),
				"explorer": "https://testnet-blockexplorer.helachain.com/tx/" + txHash.Hex(),
			},
		})
	}
}

func (h BankHandler) GetContractBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		balance, err := h.bankSvc.GetContractBalance()
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, Response{
				Success: false,
				Message: "Cannot get contract balance: " + err.Error(),
			})
			return
		}

		writeJSON(w, http.StatusOK, Response{
			Success: true,
			Message: "Get contract balance successfully",
			Data: map[string]string{
				"balance": balance.String(),
			},
		})
	}
}

func (h BankHandler) GetHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := chi.URLParam(r, "address")
		if !common.IsHexAddress(address) {
			writeJSON(w, http.StatusBadRequest, Response{
				Success: false,
				Message: "Invalid address",
			})
			return
		}

		histories, err := h.bankSvc.GetHistory(r.Context(), address)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, Response{
				Success: false,
				Message: "Cannot get history: " + err.Error(),
			})
			return
		}

		items := make([]HistoryItem, 0, len(histories))
		for _, item := range histories {
			items = append(items, HistoryItem{
				Address:   item.Address,
				Action:    item.Action,
				Amount:    item.Amount.String(),
				TxHash:    item.TXHash,
				Status:    item.Status,
				CreatedAt: item.CreatedAt,
			})
		}

		writeJSON(w, http.StatusOK, Response{
			Success: true,
			Message: "Get history successfully",
			Data: map[string]any{
				"address": address,
				"items":   items,
			},
		})
	}
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
