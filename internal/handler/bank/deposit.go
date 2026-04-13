package bank

import (
	"encoding/json"
	txpkg "hela-bank-sc/internal/blockchain/transaction"
	"log"
	"net/http"
)

type DepositRequest struct {
	Amount float64 `json:"amount"`
}

func (h Handler) Deposit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req DepositRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if !isValidAmount(req.Amount) {
			writeError(w, http.StatusBadRequest, "Amount must be greater than 0")
			return
		}

		txHash, err := h.bankSvc.Deposit(r.Context(), req.Amount)
		if err != nil {
			log.Printf("deposit failed: %v", err)
			writeError(w, http.StatusInternalServerError, "Cannot deposit")
			return
		}

		go txpkg.WaitForTx(txHash)

		writeSuccess(w, http.StatusOK, "Deposit submitted", map[string]string{
			"tx_hash":  txHash.Hex(),
			"explorer": "https://testnet-blockexplorer.helachain.com/tx/" + txHash.Hex(),
		})
	}
}
