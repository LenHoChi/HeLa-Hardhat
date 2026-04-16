package bank

import (
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/http/httptest"
	"testing"

	"hela-bank-sc/internal/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerGetContractBalance(t *testing.T) {
	type getContractBalanceMockConfig struct {
		called  bool
		balance *big.Int
		err     error
	}

	balance := big.NewInt(999999)

	tests := []struct {
		name        string
		mock        getContractBalanceMockConfig
		wantStatus  int
		wantSuccess bool
		wantMessage string
		wantData    map[string]string
	}{
		{
			name: "service error",
			mock: getContractBalanceMockConfig{
				called:  true,
				balance: nil,
				err:     errors.New("rpc failed"),
			},
			wantStatus:  http.StatusInternalServerError,
			wantSuccess: false,
			wantMessage: "Cannot get contract balance",
		},
		{
			name: "success",
			mock: getContractBalanceMockConfig{
				called:  true,
				balance: balance,
				err:     nil,
			},
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantMessage: "Get contract balance successfully",
			wantData: map[string]string{
				"balance": balance.String(),
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewService(t)
			if tc.mock.called {
				svc.On("GetContractBalance", mock.Anything).
					Return(tc.mock.balance, tc.mock.err)
			}

			handler := New(svc)
			req := httptest.NewRequest(http.MethodGet, "/contract-balance", nil)
			rec := httptest.NewRecorder()

			handler.GetContractBalance().ServeHTTP(rec, req)

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
