// +build integration

package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/landrisek/platform-go-challenge/internal/models"
	"github.com/landrisek/platform-go-challenge/internal/vault"
)

func TestControllerIntegration(t *testing.T) {

	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	serverPort := os.Getenv("ASSET_PORT")
	vaultConfig := vault.VaultConfig{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount: os.Getenv("VAULT_MOUNT"),
	}
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}
	dbConfig := models.DBConfig{
		Host:       os.Getenv("MYSQL_HOST"),
		Port:       port,
		Database:   os.Getenv("MYSQL_DATABASE"),
	}

	//assetAddr := fmt.Sprintf("http://localhost:%s", serverPort)
	blacklistAddr := fmt.Sprintf("http://localhost:%s", os.Getenv("BLACKLIST_PORT"))

	// Start the asset server with the mock DB
	go func() {
		err := RunAsset(vaultConfig, dbConfig, redisAddr, blacklistAddr, serverPort)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Wait for the server to start
	time.Sleep(10 * time.Second)

	// TODO: create test user

	// Define the test cases
	testCases := []struct {
		name               string
		method             string
		path               string
		requestBody        string
		expectedCode       int
		expectedData       string
		expectedStatusCode int
	}{
		{
			name:              "create assets",
			method:             http.MethodPost,
			path:               "/create",
			requestBody:        `[
				{
				  "id": 1,
				  "charts": [
					{
					  "title": "Chart 1",
					  "axes_titles": "X-Axis, Y-Axis",
					  "data": "1,2,3,4,5",
					  "description": "Chart 1 of user 1"
					},
					{
					  "title": "Chart 2",
					  "axes_titles": "X-Axis, Y-Axis",
					  "data": "5,4,3,2,1",
					  "description": "Chart 2 of user 1"
					}
				  ],
				  "insights": [
					{
					  "title": "Insight 1",
					  "text": "This is Insight 1",
					  "description": "Insight 1 of user 1"
					},
					{
					  "title": "Insight 2",
					  "text": "This is Insight 2",
					  "description": "Insight 2 of user 1"
					}
				  ],
				  "audiences": [
					{
					  "title": "Audience 1",
					  "characteristics": "Age: 25-35, Gender: Male",
					  "description": "This is Audience 1"
					},
					{
					  "title": "Audience 2",
					  "characteristics": "Age: 18-24, Gender: Female",
					  "description": "This is Audience 2"
					}
				  ]
				}
			  ]`,
			expectedCode:       http.StatusOK,
			expectedData:       `{"format":"","data":null}`,
			expectedStatusCode: http.StatusOK,
		},
		/*{
			method:             http.MethodPost,
			path:               "/create",
			requestBody:        `{"format": "json", "data": {"users": []}}`,
			expectedCode:       http.StatusOK,
			expectedData:       `{"format": "json", "data": []}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			method:             http.MethodPost,
			path:               "/create",
			requestBody:        `{"format": "json", "data": {"users": [{"id": 1, "name": "John Snow"}]}}`,
			expectedCode:       http.StatusUnauthorized,
			expectedData:       "",
			expectedStatusCode: http.StatusUnauthorized,
		},
		{
			method:             http.MethodPut,
			path:               "/update",
			requestBody:        `{"format": "json", "data": {"id": 1, "name": "Updated Name"}}`,
			expectedCode:       http.StatusOK,
			expectedData:       `{"format": "json", "data": "Update successful"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			method:             http.MethodPut,
			path:               "/update",
			requestBody:        `{"format": "json", "data": {"id": 2, "name": "Updated Name"}}`,
			expectedCode:       http.StatusOK,
			expectedData:       `{"format": "json", "data": "No rows updated"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			method:             http.MethodPut,
			path:               "/update",
			requestBody:        `{"format": "json", "data": {"id": 1}}`,
			expectedCode:       http.StatusBadRequest,
			expectedData:       "",
			expectedStatusCode: http.StatusBadRequest,
		},
		{
			method:             http.MethodDelete,
			path:               "/delete",
			requestBody:        `{"format": "json", "data": {"id": 1}}`,
			expectedCode:       http.StatusOK,
			expectedData:       `{"format": "json", "data": "Delete successful"}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			method:             http.MethodDelete,
			path:               "/delete",
			requestBody:        `{"format": "json", "data": {"id": 2}}`,
			expectedCode:       http.StatusOK,
			expectedData:       `{"format": "json", "data": "No rows deleted"}`,
			expectedStatusCode: http.StatusOK,
		},*/
	}

	// Iterate over the test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			// Create a request with the test case data
			requestBody := strings.NewReader(testCase.requestBody)
			request, err := http.NewRequest(testCase.method, "http://localhost:"+serverPort+testCase.path, requestBody)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Set the Authorization header with the bearer token
			request.Header.Set("Authorization", "Bearer XXX")

			// Make the request to the server
			response, err := http.DefaultClient.Do(request)
			if err != nil {
				t.Fatalf("Failed to make request: %v", err)
			}
			defer response.Body.Close()

			// Check the response status code
			assert.Equal(t, testCase.expectedStatusCode, response.StatusCode)

			// Optionally, you can read and assert the response body
			responseBody, err := ioutil.ReadAll(response.Body)
			if err != nil {
				t.Fatalf("Failed to read response body: %v", err)
			}
			assert.Equal(t, testCase.expectedData, string(responseBody))

		})
	}

	// TODO: delete test user

}