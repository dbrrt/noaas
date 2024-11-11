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

	resp, err := http.Get(fmt.Sprintf("%s", ts.URL))

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	assert.Equal(t, resp.Status, string("200 OK"))

	val, ok := resp.Header["Content-Type"]

	assert.Equal(t, ok, true)
	assert.Equal(t, val, []string([]string{"text/plain; charset=utf-8"}))
}

// func TestIntRangeToRomanConversionValid(t *testing.T) {

// 	ts := httptest.NewServer(setupServer())
// 	defer ts.Close()

// 	testRequestInput := `{"start_range": 1, "end_range": 25}`

// 	resp, err := http.Post(
// 		fmt.Sprintf("%s/v1/int/range/convert", ts.URL),
// 		"application/json", bytes.NewBuffer([]byte(testRequestInput)))

// 	assert.Equal(t, err, nil)

// 	defer resp.Body.Close()
// 	requestResponse1, err := ioutil.ReadAll(resp.Body)

// 	assert.Equal(t, err, nil)

// 	var responseObj ResponseStruct
// 	if err := json.Unmarshal(requestResponse1, &responseObj); err != nil {
// 		panic(err)
// 	}

// 	assert.Equal(t, len(responseObj.Payload), 25)
// 	assert.Equal(t, responseObj.Payload[0], "I")
// 	assert.Equal(t, len(responseObj.Error), 0)

// }
