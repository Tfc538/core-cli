//go:build e2e
// +build e2e

package main

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func coreBinaryPath(t *testing.T) string {
	t.Helper()

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get working directory: %v", err)
	}

	return filepath.Clean(filepath.Join(wd, "..", "..", "dist", "core", "core"))
}

// TestE2E_VersionCommand tests the `core version` command.
func TestE2E_VersionCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	cmd := exec.Command(coreBinaryPath(t), "version")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("core version failed: %v", err)
	}

	outputStr := string(output)
	if !contains(outputStr, "CORE CLI") {
		t.Errorf("Expected 'CORE CLI' in output, got: %s", outputStr)
	}

	if !contains(outputStr, "Commit:") {
		t.Errorf("Expected 'Commit:' in output, got: %s", outputStr)
	}

	if !contains(outputStr, "Built:") {
		t.Errorf("Expected 'Built:' in output, got: %s", outputStr)
	}
}

// TestE2E_VersionJSONOutput tests the `core version --json` command.
func TestE2E_VersionJSONOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	cmd := exec.Command(coreBinaryPath(t), "version", "--json")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("core version --json failed: %v", err)
	}

	var versionInfo struct {
		Version   string `json:"version"`
		Commit    string `json:"commit"`
		BuildDate string `json:"build_date"`
	}

	err = json.Unmarshal(output, &versionInfo)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if versionInfo.Version == "" {
		t.Error("Version field should not be empty")
	}

	if versionInfo.Commit == "" {
		t.Error("Commit field should not be empty")
	}

	if versionInfo.BuildDate == "" {
		t.Error("BuildDate field should not be empty")
	}
}

// TestE2E_UpdateCheckCommand tests the `core update check` command.
func TestE2E_UpdateCheckCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	cmd := exec.Command(coreBinaryPath(t), "update", "check")
	cmd.Env = append(os.Environ(), "CORE_UPDATE_API_BASE=http://127.0.0.1:9")
	output, err := cmd.Output()
	if err != nil {
		// It's OK if check fails (no release yet), as long as the command runs
		t.Logf("core update check returned error (expected if no releases): %v", err)
	}

	outputStr := string(output)
	// Should contain at least version info
	if !contains(outputStr, "version") && !contains(outputStr, "Error") {
		t.Logf("Unexpected output: %s", outputStr)
	}
}

// TestE2E_UpdateCheckJSONOutput tests the `core update check --json` command.
func TestE2E_UpdateCheckJSONOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	cmd := exec.Command(coreBinaryPath(t), "update", "check", "--json")
	cmd.Env = append(os.Environ(), "CORE_UPDATE_API_BASE=http://127.0.0.1:9")
	output, err := cmd.Output()

	// Parse as JSON if successful
	if err == nil && len(output) > 0 {
		var updateInfo struct {
			CurrentVersion  string `json:"current_version"`
			LatestVersion   string `json:"latest_version"`
			UpdateAvailable bool   `json:"update_available"`
		}

		err := json.Unmarshal(output, &updateInfo)
		if err != nil {
			t.Fatalf("Failed to unmarshal JSON: %v", err)
		}

		if updateInfo.CurrentVersion == "" {
			t.Error("CurrentVersion should not be empty")
		}
	}
}

// TestE2E_HelpCommand tests the `core --help` command.
func TestE2E_HelpCommand(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	cmd := exec.Command(coreBinaryPath(t), "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("core --help failed: %v", err)
	}

	outputStr := string(output)
	expectedStrings := []string{"CORE CLI", "Usage:", "Commands:"}

	for _, expected := range expectedStrings {
		if !contains(outputStr, expected) {
			t.Errorf("Expected '%s' in help output, got: %s", expected, outputStr)
		}
	}
}

// TestE2E_VersionSubcommand tests the `core version --help` command.
func TestE2E_VersionSubcommandHelp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	cmd := exec.Command(coreBinaryPath(t), "version", "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("core version --help failed: %v", err)
	}

	outputStr := string(output)
	if !contains(outputStr, "version") || !contains(outputStr, "Usage:") {
		t.Errorf("Expected version help output, got: %s", outputStr)
	}
}

// TestE2E_UpdateSubcommand tests the `core update --help` command.
func TestE2E_UpdateSubcommandHelp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	cmd := exec.Command(coreBinaryPath(t), "update", "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("core update --help failed: %v", err)
	}

	outputStr := string(output)
	expectedStrings := []string{"update", "check", "apply"}

	for _, expected := range expectedStrings {
		if !contains(outputStr, expected) {
			t.Errorf("Expected '%s' in update help, got: %s", expected, outputStr)
		}
	}
}

// TestE2E_UpdateCheckSubcommandHelp tests the `core update check --help` command.
func TestE2E_UpdateCheckSubcommandHelp(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	cmd := exec.Command(coreBinaryPath(t), "update", "check", "--help")
	output, err := cmd.Output()
	if err != nil {
		t.Fatalf("core update check --help failed: %v", err)
	}

	outputStr := string(output)
	if !contains(outputStr, "check") {
		t.Errorf("Expected 'check' in help output, got: %s", outputStr)
	}
}

// TestE2E_ExitCodes tests that commands exit with appropriate codes.
func TestE2E_ExitCodes(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	// Successful command should exit with 0
	cmd := exec.Command("./core", "version")
	err := cmd.Run()
	if err != nil {
		t.Errorf("core version should exit with 0, got error: %v", err)
	}

	// Invalid command should exit with non-zero
	cmd = exec.Command("./core", "nonexistent")
	err = cmd.Run()
	if err == nil {
		t.Error("core nonexistent should exit with non-zero code")
	}
}

// TestE2E_StdoutOutput tests that commands output to stdout correctly.
func TestE2E_StdoutOutput(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	cmd := exec.Command("./core", "version")
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		t.Fatalf("core version failed: %v", err)
	}

	if stdout.Len() == 0 {
		t.Error("Expected output on stdout")
	}

	if stderr.Len() != 0 {
		t.Logf("Unexpected output on stderr: %s", stderr.String())
	}
}

// TestE2E_CommandArgumentParsing tests various argument combinations.
func TestE2E_CommandArgumentParsing(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "version command",
			args:    []string{"version"},
			wantErr: false,
		},
		{
			name:    "version with json flag",
			args:    []string{"version", "--json"},
			wantErr: false,
		},
		{
			name:    "update check command",
			args:    []string{"update", "check"},
			wantErr: true, // May fail if no releases, but command should be valid
		},
		{
			name:    "help command",
			args:    []string{"--help"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := exec.Command("./core", tt.args...)
			err := cmd.Run()

			if tt.wantErr {
				if err == nil {
					t.Logf("Expected error for %v, but command succeeded", tt.args)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %v: %v", tt.args, err)
				}
			}
		})
	}
}

// TestE2E_BinaryExists tests that the binary exists and is executable.
func TestE2E_BinaryExists(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping e2e test in short mode")
	}

	binaryPath := "./core"
	info, err := os.Stat(binaryPath)
	if err != nil {
		t.Fatalf("core binary not found: %v", err)
	}

	if info.IsDir() {
		t.Error("core is a directory, not a binary")
	}

	// Check if executable
	if (info.Mode() & 0111) == 0 {
		t.Error("core binary is not executable")
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
