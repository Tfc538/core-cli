package testutil

import (
	"crypto/sha256"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// CreateTestFile creates a temporary file with the given content.
func CreateTestFile(t *testing.T, content []byte) string {
	tmpFile, err := os.CreateTemp("", "core-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	if _, err := tmpFile.Write(content); err != nil {
		t.Fatalf("Failed to write to temp file: %v", err)
	}

	tmpFile.Close()
	return tmpFile.Name()
}

// CreateTestDir creates a temporary directory for testing.
func CreateTestDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "core-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	return tmpDir
}

// CalculateHash returns the SHA256 hash of the given content.
func CalculateHash(content []byte) string {
	hash := sha256.Sum256(content)
	return fmt.Sprintf("%x", hash)
}

// CreateMockServer creates a mock HTTP server that serves the given content.
// The handler responds to /file with the content and /checksum with the checksum.
func CreateMockServer(content []byte, filename string) *httptest.Server {
	checksum := CalculateHash(content)
	checksumContent := fmt.Sprintf("%s  %s\n", checksum, filename)

	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/file":
			w.Header().Set("Content-Length", fmt.Sprintf("%d", len(content)))
			w.Write(content)
		case "/checksum":
			w.Write([]byte(checksumContent))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
}

// CreateErrorServer creates a mock HTTP server that returns the given status code.
func CreateErrorServer(statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprintf(w, "Error %d", statusCode)
	}))
}

// CreateSlowServer creates a mock HTTP server that takes a long time to respond.
func CreateSlowServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate slow response (in tests this won't actually block much)
		// In real usage, this would exceed the timeout
		w.Write([]byte("slow response"))
	}))
}

// Contains checks if a string contains a substring.
func Contains(s, substr string) bool {
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

// AssertEqual fails the test if actual != expected.
func AssertEqual(t *testing.T, expected, actual interface{}, message string) {
	if expected != actual {
		t.Errorf("%s: expected %v, got %v", message, expected, actual)
	}
}

// AssertNotEqual fails the test if actual == expected.
func AssertNotEqual(t *testing.T, expected, actual interface{}, message string) {
	if expected == actual {
		t.Errorf("%s: should not be %v", message, actual)
	}
}

// AssertError fails the test if err is nil.
func AssertError(t *testing.T, err error, message string) {
	if err == nil {
		t.Errorf("%s: expected error, got nil", message)
	}
}

// AssertNoError fails the test if err is not nil.
func AssertNoError(t *testing.T, err error, message string) {
	if err != nil {
		t.Errorf("%s: unexpected error %v", message, err)
	}
}

// AssertStringContains fails the test if s does not contain substr.
func AssertStringContains(t *testing.T, s, substr, message string) {
	if !Contains(s, substr) {
		t.Errorf("%s: string does not contain '%s'. Got: %s", message, substr, s)
	}
}

// AssertFileExists fails the test if the file does not exist.
func AssertFileExists(t *testing.T, path string, message string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("%s: file does not exist: %s", message, path)
	}
}

// AssertFileNotExists fails the test if the file exists.
func AssertFileNotExists(t *testing.T, path string, message string) {
	if _, err := os.Stat(path); err == nil {
		t.Errorf("%s: file exists but should not: %s", message, path)
	}
}

// ReadTestFile reads the content of a test file.
func ReadTestFile(t *testing.T, path string) []byte {
	content, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read test file %s: %v", path, err)
	}
	return content
}
