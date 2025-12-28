package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	backendversion "github.com/Tfc538/core-cli/internal/backend/service/version"
)

func TestVersionHandlerLatest(t *testing.T) {
	svc := backendversion.NewService(backendversion.NewInMemoryProvider([]backendversion.Info{{
		Version:   "0.1.0",
		Commit:    "abc123",
		BuildDate: "2025-01-01T00:00:00Z",
	}}))

	handler := VersionHandler{Service: svc}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/version/latest", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected data map, got %T", resp.Data)
	}
	if data["version"] != "0.1.0" {
		t.Fatalf("expected version 0.1.0, got %v", data["version"])
	}
}

func TestVersionHandlerSpecific(t *testing.T) {
	svc := backendversion.NewService(backendversion.NewInMemoryProvider([]backendversion.Info{{
		Version:   "0.2.0",
		Commit:    "def456",
		BuildDate: "2025-02-01T00:00:00Z",
	}}))

	handler := VersionHandler{Service: svc}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/version/0.2.0", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", rec.Code)
	}

	var resp Response
	if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	data, ok := resp.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected data map, got %T", resp.Data)
	}
	if data["version"] != "0.2.0" {
		t.Fatalf("expected version 0.2.0, got %v", data["version"])
	}
}

func TestVersionHandlerNotFound(t *testing.T) {
	svc := backendversion.NewService(backendversion.NewInMemoryProvider([]backendversion.Info{{
		Version:   "0.3.0",
		Commit:    "ghi789",
		BuildDate: "2025-03-01T00:00:00Z",
	}}))

	handler := VersionHandler{Service: svc}

	req := httptest.NewRequest(http.MethodGet, "/api/v1/version/9.9.9", nil)
	rec := httptest.NewRecorder()
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", rec.Code)
	}
}
