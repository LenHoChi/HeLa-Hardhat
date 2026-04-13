package bank

import (
	"encoding/json"
	txpkg "hela-bank-sc/internal/blockchain/transaction"
	"log"
	"net/http"
)

type WithdrawRequest struct {
	Amount float64 `json:"amount"`
}

func (h Handler) Withdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req WithdrawRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, "Invalid request body")
			return
		}

		if !isValidAmount(req.Amount) {
			writeError(w, http.StatusBadRequest, "Amount must be greater than 0")
			return
		}

		txHash, err := h.bankSvc.Withdraw(r.Context(), req.Amount)
		if err != nil {
			log.Printf("withdraw failed: %v", err)
			writeError(w, http.StatusInternalServerError, "Cannot withdraw")
			return
		}

		go txpkg.WaitForTx(txHash)

		writeSuccess(w, http.StatusOK, "Withdraw submitted", map[string]string{
			"tx_hash":  txHash.Hex(),
			"explorer": "https://testnet-blockexplorer.helachain.com/tx/" + txHash.Hex(),
		})
	}
}
