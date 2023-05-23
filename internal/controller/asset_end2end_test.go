//go:build end2end
// +build end2end

package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/landrisek/platform-go-challenge/internal/models"
	"github.com/landrisek/platform-go-challenge/internal/sagas"
	"github.com/landrisek/platform-go-challenge/internal/vault"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestAsset(t *testing.T) {

	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	serverPort := "8090"
	userPort := "9090"
	token := "XXX"
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
	userAddr := fmt.Sprintf("http://localhost:%s", userPort)
	blacklistAddr := fmt.Sprintf("http://localhost:%s", os.Getenv("BLACKLIST_PORT"))

	// Create the database connection string
	// HINT: Here we are using DB connection without vault, directly taken from env variables
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database)
	// Connect to the database
	db, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		t.Fatalf("Opening DB failed: %v", err)
	}
	defer db.Close()

	err = cleanupTables(db)
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

	// Start the asset server with the mock DB
	go func() {
		err := RunUser(ctx, vaultConfig, dbConfig, userPort)
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
		requestBody        func(*testing.T, []models.User) string
		expectedCode       int
		expectedData       func(*testing.T, *sqlx.DB, []models.User, []byte) []models.User
		expectedStatusCode int
	}{
		{
			name:               "create user",
			method:             http.MethodPost,
			token:              token,
			path:               userAddr + "/create",
			requestBody:        user,
			expectedCode:       http.StatusOK,
			expectedData:       checkUser,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "create assets",
			method:             http.MethodPost,
			token:              token,
			path:               assetAddr + "/create",
			requestBody:        create,
			expectedCode:       http.StatusOK,
			expectedData:       checkCreate,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "read assets",
			method:             http.MethodPost,
			token:              token,
			path:               assetAddr + "/read",
			requestBody:        read,
			expectedCode:       http.StatusOK,
			expectedData:       checkRead,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "succesful update",
			method:             http.MethodPut,
			token:              token,
			path:               assetAddr + "/update",
			requestBody:        update,
			expectedCode:       http.StatusOK,
			expectedData:       checkUpdate,
			expectedStatusCode: http.StatusOK,
		},
		{
			name:               "succesful delete",
			method:             http.MethodDelete,
			token:              token,
			path:               assetAddr + "/delete",
			requestBody:        delete,
			expectedCode:       http.StatusOK,
			expectedData:       checkDelete,
			expectedStatusCode: http.StatusOK,
		},
	}

	var users []models.User
	// Iterate over the test cases
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			content := testCase.requestBody(t, users)
			// Create a request with the test case data
			requestBody := strings.NewReader(content)
			request, err := http.NewRequest(testCase.method, testCase.path, requestBody)
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
			users = testCase.expectedData(t, db, users, responseBody)
		})
	}
}

func user(t *testing.T, users []models.User) string {
	return `{"id": 1,"name": "John Snow"}`
}

func checkUser(t *testing.T, db *sqlx.DB, users []models.User, response []byte) []models.User {
	var user models.User
	err := json.Unmarshal(response, &user)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON file on check user: %v", err)
	}
	return users
}

func create(t *testing.T, users []models.User) string {
	contentBody, err := ioutil.ReadFile("../../artifacts/asset/create.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}
	return string(contentBody)
}

func checkCreate(t *testing.T, db *sqlx.DB, users []models.User, response []byte) []models.User {
	return users
}

func read(t *testing.T, users []models.User) string {
	return `[{"id": 1}]`
}

func checkRead(t *testing.T, db *sqlx.DB, users []models.User, response []byte) []models.User {
	var genericReq sagas.GenericResponse
	err := json.Unmarshal(response, &genericReq)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON file on check read: %v", err)
	}
	err = json.Unmarshal(genericReq.Data, &users)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON file on check read: %v", err)
	}
	assert.Equal(t, users[0].ID, int64(1))
	return users
}

func update(t *testing.T, users []models.User) string {
	assert.NotEqual(t, "test", users[0].Audiences[0].Characteristics)
	users[0].Audiences[0].Characteristics = "test"
	requestBody, err := json.Marshal(&users)
	if err != nil {
		log.Println("Error on marshaling error response:", err)
	}
	return string(requestBody)
}

func checkUpdate(t *testing.T, db *sqlx.DB, users []models.User, response []byte) []models.User {
	var genericReq sagas.GenericResponse
	err := json.Unmarshal(response, &genericReq)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON file on check update: %v", err)
	}
	var usersResp []models.User
	err = json.Unmarshal(genericReq.Data, &usersResp)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON file on check update: %v", err)
	}
	// one audience was not altered, so asset server should report that
	assert.Equal(t, "Database error on audience", usersResp[0].Audiences[0].Error)
	audience, err := models.ReadAudience(db, users[0].ID, users[0].Audiences[0].ID)
	if err != nil {
		t.Fatalf("Failed to get audience on check update: %v", err)
	}
	assert.Equal(t, "test", audience.Characteristics)
	return users
}

func delete(t *testing.T, users []models.User) string {
	requestBody, err := json.Marshal(&users)
	if err != nil {
		log.Println("Error on marshaling error response:", err)
	}
	return string(requestBody)
}

func checkDelete(t *testing.T, db *sqlx.DB, users []models.User, response []byte) []models.User {
	_, err := models.ReadAudience(db, users[0].ID, users[0].Audiences[0].ID)
	assert.NotEqual(t, nil, err)
	assert.Equal(t, "sql: no rows in result set", err.Error())
	return users
}

// HINT: cleanup tables is run only before test, not after
// in this way is easy to check manually after failing tests
// but before running test we can rely DB will be cleaned up
// this ensure reliable behavior on tests
func cleanupTables(db *sqlx.DB) error {
	// Clean up assets table
	// DELETE ON CASCADE should cleanup all underlying
	_, err := db.Exec("DELETE FROM assets")
	if err != nil {
		return err
	}
	// Clean up users table
	_, err = db.Exec("DELETE FROM users")
	if err != nil {
		return err
	}

	return nil
}
