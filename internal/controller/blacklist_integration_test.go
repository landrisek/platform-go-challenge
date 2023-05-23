//go:build integration
// +build integration

package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/landrisek/platform-go-challenge/internal/models"

	"github.com/stretchr/testify/assert"
)

func TestBlacklist(t *testing.T) {
	redisAddr := "localhost:6379"
	blacklistPort := "9091"

	errChan := make(chan error)
	defer close(errChan)

	// Start the blacklist server
	go Blacklist(redisAddr, blacklistPort, errChan)
	time.Sleep(100 * time.Millisecond) // Allow some time for the server to start

	// Define the test cases
	testCases := []struct {
		name               string
		method             string
		path               string
		requestBody        string
		expectedStatusCode int
		expectedData       []models.User
		expectedError      string
	}{
		{
			name:               "Valid request",
			method:             http.MethodPost,
			path:               "/",
			requestBody:        `[{"name": "John Snow"}, {"name": "Ygritte"}]`,
			expectedStatusCode: http.StatusOK,
			expectedData: []models.User{
				{Name: "John Snow"},
				{Name: "Ygritte"},
			},
		},
	}

	// Iterate over the test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			// Create a request with the test case data
			requestBody := []byte(testCase.requestBody)
			request, err := http.NewRequest(testCase.method, fmt.Sprintf("http://localhost:%s%s", blacklistPort, testCase.path), ioutil.NopCloser(bytes.NewReader(requestBody)))
			if err != nil {
				t.Fatalf("Failed to create request on %s: %v", testCase.name, err)
			}

			// Make the request to the server
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatalf("Failed to make request on %s: %v", testCase.name, err)
			}
			defer response.Body.Close()

			// Check the response status code
			assert.Equal(t, testCase.expectedStatusCode, response.StatusCode)

			// Optionally, you can read and assert the response body
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("Failed to read response body on %s: %v", testCase.name, err)
			}
			var users []models.User
			err = json.Unmarshal(responseBody, &users)
			if err != nil {
				t.Fatalf("Failed to unmarshal JSON file on %s: %v", testCase.name, err)
			}
			assert.Equal(t, testCase.expectedData[0].Name, users[0].Name)
			assert.Equal(t, testCase.expectedData[1].Name, users[1].Name)
		})
	}
}
