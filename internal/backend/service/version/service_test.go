package version

import (
	"context"
	"testing"
)

func TestServiceLatestAndGet(t *testing.T) {
	provider := NewInMemoryProvider([]Info{{
		Version:   "0.1.0",
		Commit:    "abc123",
		BuildDate: "2025-01-01T00:00:00Z",
	}})

	svc := NewService(provider)

	latest, err := svc.Latest(context.Background())
	if err != nil {
		t.Fatalf("expected latest version, got error: %v", err)
	}
	if latest.Version != "0.1.0" {
		t.Fatalf("expected version 0.1.0, got %s", latest.Version)
	}

	entry, ok, err := svc.Get(context.Background(), "0.1.0")
	if err != nil {
		t.Fatalf("expected version lookup to succeed, got error: %v", err)
	}
	if !ok {
		t.Fatalf("expected version to be found")
	}
	if entry.Commit != "abc123" {
		t.Fatalf("expected commit abc123, got %s", entry.Commit)
	}

	_, ok, err = svc.Get(context.Background(), "9.9.9")
	if err != nil {
		t.Fatalf("expected missing version to return ok=false, got error: %v", err)
	}
	if ok {
		t.Fatalf("expected missing version to return ok=false")
	}
}
