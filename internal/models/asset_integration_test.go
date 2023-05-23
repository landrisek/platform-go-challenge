//go:build database
// +build database

package models

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestCreateAsset(t *testing.T) {

	// Create the database connection string
	// HINT: Here we are using DB connection without vault, directly taken from env variables
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}
	dataSourceName := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		port,
		os.Getenv("MYSQL_DATABASE"))
	// Connect to the database
	db, err := sqlx.Open("mysql", dataSourceName)
	if err != nil {
		t.Fatalf("Opening DB failed: %v", err)
	}
	defer db.Close()

	insertQuery := "INSERT INTO users (id, name) VALUES (?, ?)"
	_, err = db.Exec(insertQuery, 1, "John Snow")
	if err != nil {
		t.Fatalf("Adding test user failed: %v", err)
	}

	// Define a sample asset
	asset := Asset{
		Description: "Sample asset",
		Type:        "audiences",
		UserID:      1,
	}

	// Call the createAsset function
	assetID, err := createAsset(db, asset)
	if err != nil {
		t.Fatalf("Failed to create asset: %v", err)
	}

	// Assert that the assetID is not zero (indicating a successful insertion)
	assert.NotZero(t, assetID)

	// Optional: Retrieve the asset from the database and assert its values
	var retrievedAsset Asset
	err = db.Get(&retrievedAsset, "SELECT * FROM assets WHERE id = ?", assetID)
	if err != nil {
		t.Fatalf("Failed to retrieve asset: %v", err)
	}

	assert.Equal(t, asset.Description, retrievedAsset.Description)
	assert.Equal(t, asset.Type, retrievedAsset.Type)
	assert.Equal(t, asset.UserID, retrievedAsset.UserID)
}
