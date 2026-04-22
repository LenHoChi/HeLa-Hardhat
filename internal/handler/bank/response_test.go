package bank

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

type payload struct {
	Name string `json:"name"`
}

func TestWriteJSON(t *testing.T) {
	tests := []struct {
		name       string
		status     int
		input      payload
		wantStatus int
		wantBody   payload
	}{
		{
			name:       "writes json payload",
			status:     http.StatusCreated,
			input:      payload{Name: "hela"},
			wantStatus: http.StatusCreated,
			wantBody:   payload{Name: "hela"},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			writeJSON(rec, tc.status, tc.input)

			require.Equal(t, tc.wantStatus, rec.Code)
			require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var got payload
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)
			require.Equal(t, tc.wantBody, got)
		})
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		message      string
		wantStatus   int
		wantResponse Response
	}{
		{
			name:       "writes error response",
			status:     http.StatusBadRequest,
			message:    "bad request",
			wantStatus: http.StatusBadRequest,
			wantResponse: Response{
				Success: false,
				Message: "bad request",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			writeError(rec, tc.status, tc.message)

			require.Equal(t, tc.wantStatus, rec.Code)
			require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var got Response
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)
			require.Equal(t, tc.wantResponse, got)
		})
	}
}

func TestWriteSuccess(t *testing.T) {
	tests := []struct {
		name         string
		status       int
		message      string
		data         map[string]string
		wantStatus   int
		wantResponse Response
	}{
		{
			name:       "writes success response",
			status:     http.StatusOK,
			message:    "ok",
			data:       map[string]string{"tx_hash": "0x123"},
			wantStatus: http.StatusOK,
			wantResponse: Response{
				Success: true,
				Message: "ok",
				Data: map[string]any{
					"tx_hash": "0x123",
				},
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			writeSuccess(rec, tc.status, tc.message, tc.data)

			require.Equal(t, tc.wantStatus, rec.Code)
			require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

			var got Response
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)
			require.Equal(t, tc.wantResponse, got)
		})
	}
}
