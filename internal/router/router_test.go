package router

import (
	"encoding/json"
	"math/big"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hela-bank-sc/internal/mocks"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestRouterRoutes(t *testing.T) {
	type responseBody struct {
		Message string `json:"message"`
	}

	tests := []struct {
		name        string
		method      string
		target      string
		body        string
		wantStatus  int
		wantMessage string
	}{
		{
			name:        "get balance route",
			method:      http.MethodGet,
			target:      "/balance/invalid-address",
			wantStatus:  http.StatusBadRequest,
			wantMessage: "Invalid address",
		},
		{
			name:        "deposit route",
			method:      http.MethodPost,
			target:      "/deposit",
			body:        `{"amount":0}`,
			wantStatus:  http.StatusBadRequest,
			wantMessage: "Amount must be greater than 0",
		},
		{
			name:        "withdraw route",
			method:      http.MethodPost,
			target:      "/withdraw",
			body:        `{"amount":0}`,
			wantStatus:  http.StatusBadRequest,
			wantMessage: "Amount must be greater than 0",
		},
		{
			name:        "emergency withdraw route",
			method:      http.MethodPost,
			target:      "/emergency-withdraw",
			wantStatus:  http.StatusOK,
			wantMessage: "Emergency withdraw submitted",
		},
		{
			name:        "get contract balance route",
			method:      http.MethodGet,
			target:      "/contract-balance",
			wantStatus:  http.StatusOK,
			wantMessage: "Get contract balance successfully",
		},
		{
			name:        "check health route",
			method:      http.MethodGet,
			target:      "/check-health",
			wantStatus:  http.StatusOK,
			wantMessage: "Check health ok",
		},
		{
			name:        "get history route",
			method:      http.MethodGet,
			target:      "/history/invalid-address",
			wantStatus:  http.StatusBadRequest,
			wantMessage: "Invalid address",
		},
		{
			name:        "wrong method returns method not allowed",
			method:      http.MethodGet,
			target:      "/deposit",
			wantStatus:  http.StatusMethodNotAllowed,
			wantMessage: "",
		},
		{
			name:        "unknown route returns not found",
			method:      http.MethodGet,
			target:      "/unknown",
			wantStatus:  http.StatusNotFound,
			wantMessage: "",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewService(t)
			switch tc.target {
			case "/emergency-withdraw":
				svc.On("EmergencyWithdraw", mock.Anything).Return(common.HexToHash("0x123"), nil)
			case "/contract-balance":
				svc.On("GetContractBalance", mock.Anything).Return(big.NewInt(123), nil)
			}
			r := Router{BankSvc: svc}.Routes()
			req := httptest.NewRequest(tc.method, tc.target, strings.NewReader(tc.body))
			if tc.body != "" {
				req.Header.Set("Content-Type", "application/json")
			}
			rec := httptest.NewRecorder()

			r.ServeHTTP(rec, req)

			require.Equal(t, tc.wantStatus, rec.Code)

			if tc.wantMessage == "" {
				return
			}

			var got responseBody
			err := json.Unmarshal(rec.Body.Bytes(), &got)
			require.NoError(t, err)
			require.Equal(t, tc.wantMessage, got.Message)
		})
	}
}
