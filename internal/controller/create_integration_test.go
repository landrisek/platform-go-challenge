//go:build integration
// +build integration

package controller

import (
	"context"
	"database/sql"
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

func TestCreate(t *testing.T) {

	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	serverPort := "8091"
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

	err = cleanupTables(dbConfig, true)

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
			name:               "create assets",
			method:             http.MethodPost,
			token:              "XXX",
			path:               "/create",
			requestBody:        "create.json",
			expectedData:       `{"format":"json","data":null}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			method:             http.MethodPost,
			token:              "XXX",
			path:               "/create",
			requestBody:        "empty-create.json",
			expectedData:       `Internal Server Error`,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			method:             http.MethodPost,
			token:              "YYY",
			path:               "/create",
			requestBody:        "empty-create.json",
			expectedData:       "Unauthorized",
			expectedStatusCode: http.StatusUnauthorized,
		},
	}

	// Iterate over the test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			// Read the JSON content from the file
			contentBody, err := ioutil.ReadFile("../../artifacts/asset/" + testCase.requestBody)
			if err != nil {
				t.Fatalf("Failed to read JSON file: %v", err)
			}

			// Create a request with the test case data
			requestBody := strings.NewReader(string(contentBody))
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

func cleanupTables(dbConfig models.DBConfig, withUser bool) error {
	// Create the database connection string
	// HINT: Here we are using DB connection without vault, directly taken from env variables
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database)
	// Connect to the database
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return err
	}
	defer db.Close()

	// Clean up assets table
	// DELETE ON CASCADE should cleanup all underlying
	_, err = db.Exec("DELETE FROM assets")
	if err != nil {
		return err
	}
	// Clean up users table
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		return err
	}

	if withUser {
		// Insert user into the table
		insertQuery := "INSERT INTO users (id, name) VALUES (?, ?)"
		_, err = db.Exec(insertQuery, 1, "John Snow")
		if err != nil {
			return err
		}
	}

	return nil
}
