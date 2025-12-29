package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	backendversion "github.com/Tfc538/core-cli/internal/backend/service/version"
)

func TestIntegrationEndpoints(t *testing.T) {
	service := backendversion.NewService(backendversion.NewInMemoryProvider([]backendversion.Info{{
		Version:   "1.2.3",
		Commit:    "abc123",
		BuildDate: "2025-01-01T00:00:00Z",
	}}))

	handler := NewHandler(HandlerOptions{
		ServiceName: "core-backend",
		Version:     service,
	})

	server := httptest.NewServer(handler)
	defer server.Close()

	assertOK := func(path string) Response {
		resp, err := http.Get(server.URL + path)
		if err != nil {
			t.Fatalf("failed to GET %s: %v", path, err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected status 200 for %s, got %d", path, resp.StatusCode)
		}

		var payload Response
		if err := json.NewDecoder(resp.Body).Decode(&payload); err != nil {
			t.Fatalf("failed to decode response for %s: %v", path, err)
		}

		return payload
	}

	health := assertOK("/healthz")
	data, ok := health.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected health data map, got %T", health.Data)
	}
	if data["status"] != "ok" {
		t.Fatalf("expected health status ok, got %v", data["status"])
	}

	latest := assertOK("/api/v1/version/latest")
	latestData, ok := latest.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected version data map, got %T", latest.Data)
	}
	if latestData["version"] != "1.2.3" {
		t.Fatalf("expected version 1.2.3, got %v", latestData["version"])
	}

	specific := assertOK("/api/v1/version/1.2.3")
	specificData, ok := specific.Data.(map[string]interface{})
	if !ok {
		t.Fatalf("expected version data map, got %T", specific.Data)
	}
	if specificData["version"] != "1.2.3" {
		t.Fatalf("expected version 1.2.3, got %v", specificData["version"])
	}
}
