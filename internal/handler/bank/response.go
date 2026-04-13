package bank

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, Response{
		Success: false,
		Message: message,
	})
}

func writeSuccess(w http.ResponseWriter, status int, message string, data any) {
	writeJSON(w, status, Response{
		Success: true,
		Message: message,
		Data:    data,
	})
}
