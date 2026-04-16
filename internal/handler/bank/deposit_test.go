package bank

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"hela-bank-sc/internal/mocks"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerDeposit(t *testing.T) {
	type depositMockConfig struct {
		called bool
		amount float64
		txHash common.Hash
		err    error
	}

	txHash := common.HexToHash("0x123")

	tests := []struct {
		name         string
		body         string
		mock         depositMockConfig
		wantStatus   int
		wantSuccess  bool
		wantMessage  string
		wantTxHash   string
		wantExplorer string
	}{
		{
			name:        "invalid request body",
			body:        "{",
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Invalid request body",
		},
		{
			name:        "invalid amount",
			body:        `{"amount":0}`,
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Amount must be greater than 0",
		},
		{
			name: "service error",
			body: `{"amount":1.5}`,
			mock: depositMockConfig{
				called: true,
				amount: 1.5,
				txHash: common.Hash{},
				err:    errors.New("rpc failed"),
			},
			wantStatus:  http.StatusInternalServerError,
			wantSuccess: false,
			wantMessage: "Cannot deposit",
		},
		{
			name: "success",
			body: `{"amount":2.25}`,
			mock: depositMockConfig{
				called: true,
				amount: 2.25,
				txHash: txHash,
				err:    nil,
			},
			wantStatus:   http.StatusOK,
			wantSuccess:  true,
			wantMessage:  "Deposit submitted",
			wantTxHash:   txHash.Hex(),
			wantExplorer: "https://testnet-blockexplorer.helachain.com/tx/" + txHash.Hex(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewService(t)
			if tc.mock.called {
				svc.On("Deposit", mock.Anything, tc.mock.amount).
					Return(tc.mock.txHash, tc.mock.err)
			}

			handler := New(svc)
			req := httptest.NewRequest(http.MethodPost, "/deposit", strings.NewReader(tc.body))
			rec := httptest.NewRecorder()

			handler.Deposit().ServeHTTP(rec, req)

			assertResponse(t, rec, tc.wantStatus, tc.wantSuccess, tc.wantMessage)

			if tc.wantTxHash == "" && tc.wantExplorer == "" {
				return
			}

			var resp Response
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			require.NoError(t, err)

			data, ok := resp.Data.(map[string]any)
			require.True(t, ok, "expected response data object, got %T", resp.Data)

			require.Equal(t, tc.wantTxHash, data["tx_hash"])

			require.Equal(t, tc.wantExplorer, data["explorer"])
		})
	}
}

func assertResponse(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantSuccess bool, wantMessage string) {
	t.Helper()

	require.Equal(t, wantStatus, rec.Code)

	var resp Response
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, wantSuccess, resp.Success)

	require.Equal(t, wantMessage, resp.Message)
}
