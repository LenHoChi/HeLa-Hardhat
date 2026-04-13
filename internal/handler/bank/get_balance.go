package bank

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-chi/chi/v5"
)

func (h Handler) GetBalance() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := chi.URLParam(r, "address")
		if !isValidAddress(address) {
			writeError(w, http.StatusBadRequest, "Invalid address")
			return
		}
		userAddr := common.HexToAddress(address)
		balance, err := h.bankSvc.GetBalance(userAddr)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Cannot get balance")
			return
		}

		writeSuccess(w, http.StatusOK, "Get balance successfully", map[string]string{
			"address": address,
			"balance": balance.String(),
		})
	}
}
