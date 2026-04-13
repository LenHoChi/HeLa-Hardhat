package bank

import "net/http"

func (h Handler) CheckHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		writeSuccess(w, http.StatusOK, "Check health ok", nil)
	}
}
