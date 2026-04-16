package bank

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"hela-bank-sc/internal/mocks"
)

func TestHandlerCheckHealth(t *testing.T) {
	svc := mocks.NewService(t)
	handler := New(svc)
	req := httptest.NewRequest(http.MethodGet, "/check-health", nil)
	rec := httptest.NewRecorder()

	handler.CheckHealth().ServeHTTP(rec, req)

	assertResponse(t, rec, http.StatusOK, true, "Check health ok")
}
