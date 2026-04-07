package handler

import (
	"encoding/json"
	bank "hela-bank-sc/internal/blockchain"
	"net/http"

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

// func GetBalance(w http.ResponseWriter, r *http.Request) {
// 	// txHash, err := bank.GetBalance(req.Amount)
// 	writeJSON(w, http.StatusOK, Response{
// 		Success: true,
// 		Message: "Get balance successfully",
// 	})
// }

// func Deposit(w http.ResponseWriter, r *http.Request) {
// 	// txHash, err := bank.Deposit(req.Amount)
// 	writeJSON(w, http.StatusOK, Response{
// 		Success: true,
// 		Message: "Deposit submitted",
// 	})
// }

// func Withdraw(w http.ResponseWriter, r *http.Request) {
// 	// txHash, err := bank.Withdraw(req.Amount)
// 	writeJSON(w, http.StatusOK, Response{
// 		Success: true,
// 		Message: "Withdraw submitted",
// 	})
// }

func (h *BankHandler) GetBalance() http.HandlerFunc {
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

func (h *BankHandler) Deposit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DepositRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, Response{
				Success: false,
				Message: "Invalid request body",
			})
			return
		}

		txHash, err := h.bankSvc.Deposit(req.Amount)
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

func (h *BankHandler) Withdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req WithdrawRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeJSON(w, http.StatusBadRequest, Response{
				Success: false,
				Message: "Invalid request body",
			})
			return
		}

		txHash, err := h.bankSvc.Withdraw(req.Amount)
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

func (h *BankHandler) EmergencyWithdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txHash, err := h.bankSvc.EmergencyWithdraw()
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

func (h *BankHandler) GetContractBalance() http.HandlerFunc {
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

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
