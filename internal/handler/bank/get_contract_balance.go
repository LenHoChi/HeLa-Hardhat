package bank

import "net/http"

func (h Handler) GetContractBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		balance, err := h.bankSvc.GetContractBalance(r.Context())
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Cannot get contract balance")
			return
		}

		writeSuccess(w, http.StatusOK, "Get contract balance successfully", map[string]string{
			"balance": balance.String(),
		})
	}
}
