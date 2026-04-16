package bank

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"hela-bank-sc/internal/mocks"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerEmergencyWithdraw(t *testing.T) {
	type emergencyWithdrawMockConfig struct {
		called bool
		txHash common.Hash
		err    error
	}

	txHash := common.HexToHash("0x789")

	tests := []struct {
		name         string
		mock         emergencyWithdrawMockConfig
		wantStatus   int
		wantSuccess  bool
		wantMessage  string
		wantTxHash   string
		wantExplorer string
	}{
		{
			name: "service error",
			mock: emergencyWithdrawMockConfig{
				called: true,
				txHash: common.Hash{},
				err:    errors.New("rpc failed"),
			},
			wantStatus:  http.StatusInternalServerError,
			wantSuccess: false,
			wantMessage: "Cannot emergency withdraw",
		},
		{
			name: "success",
			mock: emergencyWithdrawMockConfig{
				called: true,
				txHash: txHash,
				err:    nil,
			},
			wantStatus:   http.StatusOK,
			wantSuccess:  true,
			wantMessage:  "Emergency withdraw submitted",
			wantTxHash:   txHash.Hex(),
			wantExplorer: "https://testnet-blockexplorer.helachain.com/tx/" + txHash.Hex(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewService(t)
			if tc.mock.called {
				svc.On("EmergencyWithdraw", mock.Anything).
					Return(tc.mock.txHash, tc.mock.err)
			}

			handler := New(svc)
			req := httptest.NewRequest(http.MethodPost, "/emergency-withdraw", nil)
			rec := httptest.NewRecorder()

			handler.EmergencyWithdraw().ServeHTTP(rec, req)

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
