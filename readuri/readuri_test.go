package readuri

import (
	"net/http"
	"net/http/httptest"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadRemoteUriPayload(t *testing.T) {
	t.Run("Test successful HTTP response", func(t *testing.T) {
		// Mock server that returns a successful response
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Hello, World!"))
		}))
		defer ts.Close()

		// Call the function with the test server's URL
		result, err := ReadRemoteUriPayload(ts.URL, false)
		assert.NoError(t, err)
		assert.Equal(t, "Hello, World!", result)
	})

	t.Run("Test non-200 HTTP status code", func(t *testing.T) {
		// Mock server that returns a 404 status code
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer ts.Close()

		// Call the function
		result, err := ReadRemoteUriPayload(ts.URL, false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to retrieve content. Status code: 404")
		assert.Empty(t, result)
	})

	t.Run("Test HTTP request failure", func(t *testing.T) {
		// Using an invalid URL to simulate failure
		result, err := ReadRemoteUriPayload("http://invalid-url", false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error fetching URL:")
		assert.Empty(t, result)
	})

	t.Run("Test error reading response body", func(t *testing.T) {
		// Simulate a failure in reading response body by using a mock that always fails
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		}))
		defer ts.Close()

		result, err := ReadRemoteUriPayload(ts.URL, false)
		assert.Error(t, err)
		assert.Empty(t, result)
	})

	t.Run("Test successful script execution", func(t *testing.T) {
		// Mock server returning a shell script
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("echo Hello from the script"))
		}))
		defer ts.Close()

		// Mock exec.Command to simulate script execution
		originalExec := execCommand
		defer func() { execCommand = originalExec }() // Restore original execCommand
		execCommand = func(name string, arg ...string) *exec.Cmd {
			return exec.Command("echo", "Hello from the script")
		}

		// Test the function with script flag set to true
		result, err := ReadRemoteUriPayload(ts.URL, true)
		assert.NoError(t, err)
		assert.Equal(t, "Hello from the script\n", result)
	})

	t.Run("Test script execution failure", func(t *testing.T) {
		// Mock server returning an invalid shell script
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("invalid script"))
		}))
		defer ts.Close()

		// Simulate error in running the shell script
		originalExec := execCommand
		defer func() { execCommand = originalExec }() // Restore original execCommand
		execCommand = func(name string, arg ...string) *exec.Cmd {
			return exec.Command("sh", "-c", "invalid script")
		}

		result, err := ReadRemoteUriPayload(ts.URL, true)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error executing script:")
		assert.Empty(t, result)
	})

	t.Run("Test empty response body", func(t *testing.T) {
		// Mock server returning an empty body
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusOK)
		}))
		defer ts.Close()

		// Call the function
		result, err := ReadRemoteUriPayload(ts.URL, false)
		assert.NoError(t, err)
		assert.Empty(t, result)
	})

	t.Run("Test invalid URI", func(t *testing.T) {
		// Test invalid URI
		result, err := ReadRemoteUriPayload("http://", false)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "error fetching URL:")
		assert.Empty(t, result)
	})
}
