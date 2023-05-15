package models

import (
	"fmt"

	"github.com/landrisek/platform-go-challenge/internal/vault"
	"github.com/pkg/errors"
)

type DBConfig struct {
	Host       string
	Port       int
	Database   string
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