package bank

import "net/http"

func (h Handler) EmergencyWithdraw() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		txHash, err := h.bankSvc.EmergencyWithdraw(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Cannot emergency withdraw")
			return
		}

		writeSuccess(w, http.StatusOK, "Emergency withdraw submitted", map[string]string{
			"tx_hash":  txHash.Hex(),
			"explorer": "https://testnet-blockexplorer.helachain.com/tx/" + txHash.Hex(),
		})
	}
}
