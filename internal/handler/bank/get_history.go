package bank

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
)

type HistoryItem struct {
	Address   string    `json:"address"`
	Action    string    `json:"action"`
	Amount    string    `json:"amount"`
	TxHash    string    `json:"tx_hash"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (h Handler) GetHistory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		address := chi.URLParam(r, "address")
		if !isValidAddress(address) {
			writeError(w, http.StatusBadRequest, "Invalid address")
			return
		}

		histories, err := h.bankSvc.GetHistory(r.Context(), address)
		if err != nil {
			writeError(w, http.StatusInternalServerError, "Cannot get history")
			return
		}

		items := make([]HistoryItem, 0, len(histories))
		for _, item := range histories {
			items = append(items, HistoryItem{
				Address:   item.Address,
				Action:    item.Action,
				Amount:    item.Amount,
				TxHash:    item.TxHash,
				Status:    item.Status,
				CreatedAt: item.CreatedAt,
			})
		}

		writeSuccess(w, http.StatusOK, "Get history successfully", map[string]any{
			"address": address,
			"items":   items,
		})
	}
}
