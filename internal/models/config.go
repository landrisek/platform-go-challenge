package models

import (
	"fmt"

	"github.com/landrisek/platform-go-challenge/internal/vault"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type DBConfig struct {
	Host     string
	Port     int
	Database string
}

func GetDatabaseURL(vaultConfig vault.VaultConfig, dbConfig DBConfig) (string, error) {
	// Retrieve MySQL credentials from Vault
	creds, err := vault.GetSQLCredentials(vaultConfig)
	if err != nil {
		return "", errors.Wrap(err, "failed to get credentials from Vault")
	}
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		creds["username"],
		creds["password"],
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Database), nil
}

func OpenDB(vaultConfig vault.VaultConfig, dbConfig DBConfig) (*sqlx.DB, error) {
	// Define the MySQL database connection string

	dbURL, err := GetDatabaseURL(vaultConfig, dbConfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed on vault in server")
	}
	// Open a connection to the database
	db, err := sqlx.Open("mysql", dbURL)
	if err != nil {
		return nil, errors.Wrap(err, "failed on mysql in server")
	}
	// Ping the database to check the connection
	err = db.Ping()
	if err != nil {
		return nil, errors.Wrap(err, "failed on mysql connection in server")
	}

	return db, nil
}
