package bank

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
)

func assertResponse(t *testing.T, rec *httptest.ResponseRecorder, wantStatus int, wantSuccess bool, wantMessage string) {
	t.Helper()

	require.Equal(t, wantStatus, rec.Code)

	var resp Response
	err := json.Unmarshal(rec.Body.Bytes(), &resp)
	require.NoError(t, err)

	require.Equal(t, wantSuccess, resp.Success)
	require.Equal(t, wantMessage, resp.Message)
}

func withURLParam(req *http.Request, key string, value string) *http.Request {
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(key, value)

	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
}
