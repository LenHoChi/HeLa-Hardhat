package bank

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"hela-bank-sc/internal/mocks"
)

func TestHandlerCheckHealth(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		target      string
		wantStatus  int
		wantSuccess bool
		wantMessage string
	}{
		{
			name:        "success",
			method:      http.MethodGet,
			target:      "/check-health",
			wantStatus:  http.StatusOK,
			wantSuccess: true,
			wantMessage: "Check health ok",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			svc := mocks.NewService(t)
			handler := New(svc)
			req := httptest.NewRequest(tc.method, tc.target, nil)
			rec := httptest.NewRecorder()

			handler.CheckHealth().ServeHTTP(rec, req)

			assertResponse(t, rec, tc.wantStatus, tc.wantSuccess, tc.wantMessage)
		})
	}
}
