package bank

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"hela-bank-sc/internal/mocks"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerGetBalance(t *testing.T) {
	type getBalanceMockConfig struct {
		called  bool
		addr    common.Address
		balance *big.Int
		err     error
	}

	address := "0x1234567890123456789012345678901234567890"
	balance := big.NewInt(123456789)

	tests := []struct {
		name        string
		address     string
		mock        getBalanceMockConfig
		wantStatus  int
		wantSuccess bool
		wantMessage string
		wantData    map[string]string
	}{
		{
			name:        "invalid address",
			address:     "invalid-address",
			wantStatus:  http.StatusBadRequest,
			wantSuccess: false,
			wantMessage: "Invalid address",
		},
		{
			name:    "service error",
			address: address,
			mock: getBalanceMockConfig{
				called:  true,
				addr:    common.HexToAddress(address),
				balance: nil,
				err:     errors.New("rpc failed"),
			},
			wantStatus:  http.StatusInternalServerError,
			wantSuccess: false,
			wantMessage: "Cannot get balance",
		},
		{
			name:    "success",
			address: address,
			mock: getBalanceMockConfig{
				called:  true,
				addr:    common.HexToAddress(address),
				balance: balance,
				err:     nil,
			},
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantMessage: "Get balance successfully",
			wantData: map[string]string{
				"address": address,
				"balance": balance.String(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewService(t)
			if tc.mock.called {
				svc.On("GetBalance", mock.Anything, tc.mock.addr).
					Return(tc.mock.balance, tc.mock.err)
			}

			handler := New(svc)
			req := httptest.NewRequest(http.MethodGet, "/balance/"+tc.address, nil)
			req = withURLParam(req, "address", tc.address)
			rec := httptest.NewRecorder()

			handler.GetBalance().ServeHTTP(rec, req)

			assertResponse(t, rec, tc.wantStatus, tc.wantSuccess, tc.wantMessage)

			if tc.wantData == nil {
				return
			}

			var resp Response
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			require.NoError(t, err)

			data, ok := resp.Data.(map[string]any)
			require.True(t, ok, "expected response data object, got %T", resp.Data)

			for key, want := range tc.wantData {
				require.Equal(t, want, data[key], "unexpected value for key %s", key)
			}
		})
	}
}
