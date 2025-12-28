package version

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name string
		want struct {
			version   string
			commit    string
			buildDate string
		}
	}{
		{
			name: "returns valid info",
			want: struct {
				version   string
				commit    string
				buildDate string
			}{
				version:   "dev",
				commit:    "unknown",
				buildDate: "unknown",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			info := Get()

			if info.Version == "" {
				t.Error("Version should not be empty")
			}

			if info.GitCommit == "" {
				t.Error("GitCommit should not be empty")
			}

			if info.BuildDate == "" {
				t.Error("BuildDate should not be empty")
			}
		})
	}
}

func TestInfo_String(t *testing.T) {
	tests := []struct {
		name     string
		info     Info
		contains []string
	}{
		{
			name: "formats version correctly",
			info: Info{
				Version:   "1.0.0",
				GitCommit: "abc123",
				BuildDate: "2025-12-28T10:30:00Z",
			},
			contains: []string{"CORE CLI", "v1.0.0", "Commit", "abc123", "Built", "2025-12-28"},
		},
		{
			name: "handles development version",
			info: Info{
				Version:   "dev",
				GitCommit: "unknown",
				BuildDate: "unknown",
			},
			contains: []string{"CORE CLI", "vdev", "unknown"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			str := tt.info.String()
			if str == "" {
				t.Error("String() should not return empty string")
			}

			for _, s := range tt.contains {
				if !contains(str, s) {
					t.Errorf("String should contain '%s', got: %s", s, str)
				}
			}
		})
	}
}

func TestInfo_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		info    Info
		wantErr bool
	}{
		{
			name: "marshals to valid JSON",
			info: Info{
				Version:   "1.0.0",
				GitCommit: "abc123",
				BuildDate: "2025-12-28T10:30:00Z",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := json.Marshal(tt.info)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal error = %v, wantErr %v", err, tt.wantErr)
			}

			if err == nil {
				var decoded Info
				if err := json.Unmarshal(b, &decoded); err != nil {
					t.Errorf("Unmarshal failed: %v", err)
				}

				if decoded.Version != tt.info.Version {
					t.Errorf("Version mismatch: got %s, want %s", decoded.Version, tt.info.Version)
				}
			}
		})
	}
}

func TestInfo_FieldNames(t *testing.T) {
	// Verify JSON field names match expected format
	info := Info{
		Version:   "1.0.0",
		GitCommit: "abc123",
		BuildDate: "2025-12-28T10:30:00Z",
	}

	b, _ := json.Marshal(info)
	jsonStr := string(b)

	expectedFields := map[string]bool{
		"version":    true,
		"commit":     true,
		"build_date": true,
	}

	for field := range expectedFields {
		if !strings.Contains(jsonStr, field) {
			t.Errorf("Expected JSON field '%s' not found in: %s", field, jsonStr)
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
