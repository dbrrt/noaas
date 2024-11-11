package main

import (
	"bytes"
	"dbrrt/noaas/routing"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Response structure to match the expected response from the ServiceProvisionner endpoint.
type ResponseStruct struct {
	Url   *string `json:"url"`
	Error *string `json:"error"`
}

// TestHealthRoute checks if the health route is reachable and returns the expected status and content type.
func TestHealthRoute(t *testing.T) {
	// Start a new test server with the routing setup.
	ts := httptest.NewServer(routing.SetupServer())
	defer ts.Close()

	// Send a GET request to the health route.
	resp, err := http.Get(fmt.Sprintf("%s/", ts.URL))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	// Verify the response status and content type.
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	val, ok := resp.Header["Content-Type"]
	assert.True(t, ok)
	assert.Equal(t, []string{"text/plain; charset=utf-8"}, val)
}

// TestServiceProvisionner tests various scenarios for the ServiceProvisionner endpoint.
func TestServiceProvisionner(t *testing.T) {
	// Start a new test server with the routing setup.
	ts := httptest.NewServer(routing.SetupServer())
	defer ts.Close()

	// Helper function to send a PUT request to the service route.
	sendPutRequest := func(endpoint string, body map[string]interface{}) (*http.Response, ResponseStruct) {
		reqBody, _ := json.Marshal(body)

		req, _ := http.NewRequest(http.MethodPut, fmt.Sprintf("%s%s", ts.URL, endpoint), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			t.Fatalf("Failed to send request: %v", err)
		}
		defer resp.Body.Close()

		var response ResponseStruct
		json.NewDecoder(resp.Body).Decode(&response)
		return resp, response
	}

	t.Run("Valid Request", func(t *testing.T) {
		resp, response := sendPutRequest("/v1/services/testName", map[string]interface{}{
			"url":    "http://valid-url.com",
			"script": true,
		})

		assert.Equal(t, http.StatusOK, resp.StatusCode)
		assert.NotNil(t, response.Url)
		assert.Nil(t, response.Error)
	})

	t.Run("Wrong path", func(t *testing.T) {
		resp, response := sendPutRequest("/v1/services/", map[string]interface{}{
			"url":    "http://valid-url.com",
			"script": true,
		})

		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		assert.Nil(t, response.Error)
	})

	t.Run("Invalid URI", func(t *testing.T) {
		resp, response := sendPutRequest("/v1/services/testName", map[string]interface{}{
			"url":    "not-a-valid-url",
			"script": true,
		})

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Nil(t, response.Url)
		assert.NotNil(t, response.Error)
		assert.Contains(t, *response.Error, "url")
	})

	t.Run("Missing Script Parameter", func(t *testing.T) {
		resp, response := sendPutRequest("/v1/services/testName", map[string]interface{}{
			"url": "http://valid-url.com",
		})

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Nil(t, response.Url)
		assert.NotNil(t, response.Error)
	})

	t.Run("Empty Request Body", func(t *testing.T) {
		resp, response := sendPutRequest("/v1/services/testName", nil)

		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		assert.Nil(t, response.Url)
		assert.NotNil(t, response.Error)
	})
}
