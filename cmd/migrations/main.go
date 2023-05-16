package main

import (
	"database/sql"
	"log"
	"strconv"
	"os"

	"github.com/landrisek/platform-go-challenge/internal/vault"
	"github.com/landrisek/platform-go-challenge/internal/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pressly/goose"
	"github.com/pkg/errors"
)


func main() {
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}
	err = migrate(vault.VaultConfig{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount: os.Getenv("VAULT_MOUNT"),
	},
	models.DBConfig{
		Host:       os.Getenv("MYSQL_HOST"),
		Port:       port,
		Database:   os.Getenv("MYSQL_DATABASE"),
	},
	os.Getenv("MIGRATION_DIR"))
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

func migrate(vaultConfig vault.VaultConfig, dbConfig models.DBConfig, migrateDir string) error {
	dbURL, err := models.GetDatabaseURL(vaultConfig, dbConfig)
	if err != nil {
		return errors.Wrap(err, "failed to get credentials from vault")
	}
	// Create a database connection
	db, err := sql.Open("mysql", dbURL)
	if err != nil {
		return errors.Wrap(err, "failed to connect to the database")
	}
	defer db.Close()

	// Set the MySQL dialect for Goose
	goose.SetDialect("mysql")

	// Run the migration
	err = goose.Up(db, migrateDir)
	if err != nil {
		return errors.Wrap(err, "migration failed")
	}

	return nil
}
