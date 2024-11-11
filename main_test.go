package main

import (
	"dbrrt/noaas/routing"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type ResponseStruct struct {
	Payload []string `json:"payload"`
	Error   string   `json:"error"`
}

func TestHealthRoute(t *testing.T) {
	ts := httptest.NewServer(routing.SetupServer())
	defer ts.Close()

	resp, err := http.Get(fmt.Sprintf(ts.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, resp.Status, string("200 OK"))

	val, ok := resp.Header["Content-Type"]

	assert.Equal(t, ok, true)
	assert.Equal(t, val, []string([]string{"text/plain; charset=utf-8"}))
}
