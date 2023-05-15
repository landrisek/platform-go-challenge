package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"os"

	"github.com/landrisek/platform-go-challenge/internal/vault"

	_ "github.com/go-sql-driver/mysql" // MySQL driver
	"github.com/pressly/goose"
	"github.com/pkg/errors"
)

type dbConfig struct {
	Host       string
	Port       int
	Database   string
	VaultMount string
}

func main() {
	makeRequest()
}

type VaultResponse struct {
	Data  VaultData `json:"data"`
	Lease int       `json:"lease_duration"`
}

type VaultData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

//curl --header "X-Vault-Token: myroot" --request GET  http://localhost:8200/v1/database/creds/sudo
func makeRequest() {
	fmt.Println("-------makeRequest()-------")
	req, err := http.NewRequest(http.MethodGet, "http://vault:8200/v1/mysql_sandbox/creds/sudo", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("X-Vault-Token", "myroot")
	response, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	var vaultResponse VaultResponse
	err = json.NewDecoder(response.Body).Decode(&vaultResponse)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Username:", vaultResponse.Data.Username)
	fmt.Println("Password:", vaultResponse.Data.Password)
	fmt.Println("Lease Duration:", vaultResponse.Lease)
}




func main2() {
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}
	err = migrate(vault.VaultConfig{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
	},
	dbConfig{
		Host:       os.Getenv("MYSQL_HOST"),
		Port:       port,
		Database:   os.Getenv("MYSQL_DATABASE"),
		VaultMount: os.Getenv("VAULT_MOUNT"),
	},
	os.Getenv("MIGRATION_DIR"))
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

func migrate(vaultConfig vault.VaultConfig, conf dbConfig, migrateDir string) error {

	// Retrieve MySQL credentials from Vault
	creds, err := vault.GetSQLCredentials(vaultConfig, conf.VaultMount)
	if err != nil {
		return errors.Wrap(err, "failed to get credentials from Vault")
	}

	// Set MySQL credentials
	databaseURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		creds["username"], creds["password"], conf.Host, conf.Port, conf.Database)

	// Create a database connection
	db, err := sql.Open("mysql", databaseURL)
	if err != nil {
		return errors.Wrap(err, "failed to connect to the database")
	}
	defer db.Close()

	// Run the migration
	err = goose.Up(db, migrateDir)
	if err != nil {
		return errors.Wrap(err, "migration failed")
	}

	return nil
}
