package bank

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"hela-bank-sc/internal/domain"
	"hela-bank-sc/internal/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestHandlerGetHistory(t *testing.T) {
	type getHistoryMockConfig struct {
		called  bool
		address string
		history []*domain.History
		err     error
	}

	address := "0x1234567890123456789012345678901234567890"
	createdAt := time.Date(2026, 4, 16, 10, 0, 0, 0, time.UTC)

	tests := []struct {
		name         string
		address      string
		mock         getHistoryMockConfig
		wantStatus   int
		wantSuccess  bool
		wantMessage  string
		wantAddress  string
		wantItemSize int
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
			mock: getHistoryMockConfig{
				called:  true,
				address: address,
				history: nil,
				err:     errors.New("db failed"),
			},
			wantStatus:  http.StatusInternalServerError,
			wantSuccess: false,
			wantMessage: "Cannot get history",
		},
		{
			name:    "success",
			address: address,
			mock: getHistoryMockConfig{
				called:  true,
				address: address,
				history: []*domain.History{
					{
						Address:   address,
						Action:    "deposit",
						Amount:    "1000000000000000000",
						TxHash:    "0xabc",
						Status:    "submitted",
						CreatedAt: createdAt,
					},
				},
				err: nil,
			},
			wantStatus:   http.StatusOK,
			wantSuccess:  true,
			wantMessage:  "Get history successfully",
			wantAddress:  address,
			wantItemSize: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewService(t)
			if tc.mock.called {
				svc.On("GetHistory", mock.Anything, tc.mock.address).
					Return(tc.mock.history, tc.mock.err)
			}

			handler := New(svc)
			req := httptest.NewRequest(http.MethodGet, "/history/"+tc.address, nil)
			req = withURLParam(req, "address", tc.address)
			rec := httptest.NewRecorder()

			handler.GetHistory().ServeHTTP(rec, req)

			assertResponse(t, rec, tc.wantStatus, tc.wantSuccess, tc.wantMessage)

			if tc.wantAddress == "" && tc.wantItemSize == 0 {
				return
			}

			var resp Response
			err := json.Unmarshal(rec.Body.Bytes(), &resp)
			require.NoError(t, err)

			data, ok := resp.Data.(map[string]any)
			require.True(t, ok, "expected response data object, got %T", resp.Data)

			require.Equal(t, tc.wantAddress, data["address"])

			items, ok := data["items"].([]any)
			require.True(t, ok, "expected items array, got %T", data["items"])

			require.Len(t, items, tc.wantItemSize)
		})
	}
}
