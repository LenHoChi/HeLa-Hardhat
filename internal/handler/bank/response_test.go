package bank

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteJSON(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	rec := httptest.NewRecorder()

	writeJSON(rec, http.StatusCreated, payload{Name: "hela"})

	require.Equal(t, http.StatusCreated, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var got payload
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	require.NoError(t, err)
	require.Equal(t, payload{Name: "hela"}, got)
}

func TestWriteError(t *testing.T) {
	rec := httptest.NewRecorder()

	writeError(rec, http.StatusBadRequest, "bad request")

	require.Equal(t, http.StatusBadRequest, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var got Response
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	require.NoError(t, err)
	require.Equal(t, Response{
		Success: false,
		Message: "bad request",
	}, got)
}

func TestWriteSuccess(t *testing.T) {
	rec := httptest.NewRecorder()

	writeSuccess(rec, http.StatusOK, "ok", map[string]string{
		"tx_hash": "0x123",
	})

	require.Equal(t, http.StatusOK, rec.Code)
	require.Equal(t, "application/json", rec.Header().Get("Content-Type"))

	var got Response
	err := json.Unmarshal(rec.Body.Bytes(), &got)
	require.NoError(t, err)
	require.True(t, got.Success)
	require.Equal(t, "ok", got.Message)
	require.Equal(t, map[string]any{
		"tx_hash": "0x123",
	}, got.Data)
}
