//go:build integration
// +build integration

package controller

import (
	"context"
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

	_ "github.com/go-sql-driver/mysql"
	"github.com/landrisek/platform-go-challenge/internal/models"
	"github.com/landrisek/platform-go-challenge/internal/vault"
)

func TestDelete(t *testing.T) {

	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	serverPort := "8092"
	vaultConfig := vault.VaultConfig{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount:   os.Getenv("VAULT_MOUNT"),
	}
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}
	dbConfig := models.DBConfig{
		Host:     os.Getenv("MYSQL_HOST"),
		Port:     port,
		Database: os.Getenv("MYSQL_DATABASE"),
	}

	assetAddr := fmt.Sprintf("http://localhost:%s", serverPort)
	blacklistAddr := fmt.Sprintf("http://localhost:%s", os.Getenv("BLACKLIST_PORT"))

	requestBody, _, _, err := cleanupTablesWithResponses(dbConfig)
	if err != nil {
		t.Fatalf("Cleanup tables failed: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start the asset server with the mock DB
	go func() {
		err := RunAsset(ctx, vaultConfig, dbConfig, redisAddr, blacklistAddr, serverPort)
		if err != nil {
			t.Fatal(err)
		}
	}()

	// Wait for the server to start
	time.Sleep(5 * time.Second)

	// Define the test cases
	testCases := []struct {
		name               string
		method             string
		token              string
		path               string
		requestBody        string
		expectedData       string
		expectedStatusCode int
	}{
		{
			name:               "delete assets",
			method:             http.MethodDelete,
			token:              "XXX",
			path:               "/delete",
			requestBody:        requestBody,
			expectedData:       `{"format":"json","data":null}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "delete empty assets",
			method:             http.MethodDelete,
			token:              "XXX",
			path:               "/delete",
			requestBody:        `[{"format":"json","data":{}}]`,
			expectedData:       `{"format":"json","data":null}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "delete wiht incorrect fortmat assets",
			method:             http.MethodDelete,
			token:              "XXX",
			path:               "/delete",
			requestBody:        `{"format":"json","data":{}}`,
			expectedData:       `Internal Server Error`,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "delete with incorrect token",
			method:             http.MethodDelete,
			token:              "YYY",
			path:               "/delete",
			requestBody:        `[{"format":"json","data":{}}]`,
			expectedData:       "Unauthorized",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}
	// Iterate over the test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			// Create a request with the test case data
			requestBody := strings.NewReader(testCase.requestBody)
			request, err := http.NewRequest(testCase.method, assetAddr+testCase.path, requestBody)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Set the Authorization header with the bearer token
			request.Header.Set("Authorization", fmt.Sprintf("Bearer %s", testCase.token))

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
			assert.Equal(t, testCase.expectedData, strings.TrimSpace(string(responseBody)))

		})
	}
}
