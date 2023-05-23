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

func TestUpdate(t *testing.T) {

	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	serverPort := "8094"
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

	requestBody, _, errorResponseBody, err := cleanupTablesWithResponses(dbConfig)
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
			name:               "succesful update",
			method:             http.MethodPut,
			token:              "XXX",
			path:               "/update",
			requestBody:        requestBody,
			expectedData:       `{"format":"json","data":null}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "unsuccesful update",
			method:             http.MethodPut,
			token:              "XXX",
			path:               "/update",
			requestBody:        `{"format": "json", "data": {"id": 2, "name": "Updated Name"}}`,
			expectedData:       `Internal Server Error`,
			expectedStatusCode: http.StatusInternalServerError,
		},
		{
			name:               "no assets update",
			method:             http.MethodPut,
			token:              "XXX",
			path:               "/update",
			requestBody:        `[{"format":"json","data":{"id":2,"name":"Updated Name"}}]`,
			expectedData:       `{"format":"json","data":null}`,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "duplicite update",
			method:             http.MethodPut,
			token:              "XXX",
			path:               "/update",
			requestBody:        requestBody,
			expectedData:       errorResponseBody,
			expectedStatusCode: http.StatusOK,
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

func cleanupTablesWithResponses(dbConfig models.DBConfig) (string, string, string, error) {
	// Create the database connection string
	// HINT: Here we are using DB connection without vault, directly taken from env variables
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), dbConfig.Host, dbConfig.Port, dbConfig.Database)

	// Connect to the database
	db, err := sql.Open("mysql", dataSourceName)
	if err != nil {
		return "", "", "", err
	}
	defer db.Close()

	// Insert a user
	res, err := db.Exec("INSERT INTO users (name) VALUES (?)", "John Snow")
	if err != nil {
		return "", "", "", err
	}

	// Retrieve the user ID
	userID, err := res.LastInsertId()
	if err != nil {
		return "", "", "", err
	}

	// Insert a chart
	res, err = db.Exec("INSERT INTO assets (user_id, type, description) VALUES (?, 'charts', 'Chart 1A')", userID)
	if err != nil {
		return "", "", "", err
	}

	assetID, err := res.LastInsertId()
	if err != nil {
		return "", "", "", err
	}

	// Insert chart data
	res, err = db.Exec("INSERT INTO charts (assets_id, title, axes_titles, data) VALUES (?, ?, ?, ?)", assetID, "Chart 1A", "", "")
	if err != nil {
		return "", "", "", err
	}

	// Retrieve the chart ID
	chartID, err := res.LastInsertId()
	if err != nil {
		return "", "", "", err
	}

	// Insert an insight
	res, err = db.Exec("INSERT INTO assets (user_id, type, description) VALUES (?, 'insights', 'Insight 1A')", userID)
	if err != nil {
		return "", "", "", err
	}

	assetID, err = res.LastInsertId()
	if err != nil {
		return "", "", "", err
	}

	// Insert insight data
	res, err = db.Exec("INSERT INTO insights (assets_id, text) VALUES (?, ?)", assetID, "Insight 1A")
	if err != nil {
		return "", "", "", err
	}

	// Retrieve the insight ID
	insightID, err := res.LastInsertId()
	if err != nil {
		return "", "", "", err
	}

	// Insert an audience
	res, err = db.Exec("INSERT INTO assets (user_id, type, description) VALUES (?, 'audiences', 'Audience 1A')", userID)
	if err != nil {
		return "", "", "", err
	}

	assetID, err = res.LastInsertId()
	if err != nil {
		return "", "", "", err
	}

	// Insert audience data
	res, err = db.Exec("INSERT INTO audiences (assets_id, characteristics) VALUES (?, ?)", assetID, "Audience 1A")
	if err != nil {
		return "", "", "", err
	}

	// Retrieve the audience ID
	audienceID, err := res.LastInsertId()
	if err != nil {
		return "", "", "", err
	}

	// Construct the JSON body
	request := fmt.Sprintf(`[
		{
			"id": %d,
			"charts": [
				{
					"id": %d,
					"title": "Chart 1B",
					"description": "test",
					"data": ""
				}
			],
			"insights": [
				{
					"id": %d,
					"text": "Insight 1B"
				}
			],
			"audiences": [
				{
					"id": %d,
					"characteristics": "Audience 1B"
				}
			]
		}
	]`, userID, chartID, insightID, audienceID)

	response := fmt.Sprintf(`{"format":"json","data":[{"id":%d,"name":"","audiences":[{"id":%d,"characteristics":"Audience 1A","description":"Audience 1A","error":""}],"charts":[{"id":%d,"description":"Chart 1A","error":"","AssetID":0}],"insights":[{"id":%d,"text":"Insight 1A","description":"Insight 1A","error":""}]}]}`, userID, audienceID, chartID, insightID)

	errorResponse := fmt.Sprintf(`{"format":"json","data":[{"id":%d,"name":"","audiences":[{"id":%d,"characteristics":"Audience 1B","error":"Database error on audience"}],"charts":[{"id":%d,"title":"Chart 1B","data":"","description":"test","error":"Database error on chart","AssetID":0}],"insights":[{"id":%d,"text":"Insight 1B","description":null,"error":"Database error on insight"}]}]}`, userID, chartID, insightID, audienceID)

	return request, string(response), errorResponse, nil
}
