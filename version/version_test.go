package version

import (
	"testing"
)

func TestGet(t *testing.T) {
	info := Get()

	if info.Version == "" {
		t.Error("Version should not be empty")
	}

	if info.GitCommit == "" {
		t.Error("GitCommit should not be empty (at least 'unknown')")
	}

	if info.BuildDate == "" {
		t.Error("BuildDate should not be empty (at least 'unknown')")
	}
}

func TestString(t *testing.T) {
	info := Info{
		Version:   "1.0.0",
		GitCommit: "abc123",
		BuildDate: "2025-12-28T10:30:00Z",
	}

	str := info.String()
	if str == "" {
		t.Error("String() should not return empty string")
	}

	// Verify expected content in the string
	expected := []string{"CORE CLI", "v1.0.0", "Commit", "abc123", "Built"}
	for _, s := range expected {
		if !contains(str, s) {
			t.Errorf("String should contain '%s', got: %s", s, str)
		}
	}
}

func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
