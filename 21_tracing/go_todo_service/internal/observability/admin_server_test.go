package observability

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAdminServerRoutes(t *testing.T) {
	server := NewAdminServer("127.0.0.1:0", "secret", NewMetrics(nil))

	metricsRequest := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	metricsResponse := httptest.NewRecorder()
	server.Handler.ServeHTTP(metricsResponse, metricsRequest)
	if metricsResponse.Code != http.StatusOK {
		t.Fatalf("metrics status = %d, want %d", metricsResponse.Code, http.StatusOK)
	}

	pprofRequest := httptest.NewRequest(http.MethodGet, "/debug/pprof/heap", nil)
	pprofResponse := httptest.NewRecorder()
	server.Handler.ServeHTTP(pprofResponse, pprofRequest)
	if pprofResponse.Code != http.StatusUnauthorized {
		t.Fatalf("pprof status = %d, want %d", pprofResponse.Code, http.StatusUnauthorized)
	}

	pprofRequest.Header.Set("Authorization", "Bearer secret")
	pprofResponse = httptest.NewRecorder()
	server.Handler.ServeHTTP(pprofResponse, pprofRequest)
	if pprofResponse.Code != http.StatusOK {
		t.Fatalf("authorized pprof status = %d, want %d", pprofResponse.Code, http.StatusOK)
	}
}
