package main

import (
	"fmt"
	"log"
	"os"

	"github.com/landrisek/platform-go-challenge/internal/vault"

	"github.com/pressly/goose"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

// AppConfig is used as configuration for the application
type AppConfig struct {
	Database struct {
		Host       string
		Port       int
		Database   string
		VaultMount string `yaml:"vault_mount"`
	}
}

func main() {
	err := migrate(vault.VaultConfig{
		Address: os.Getenv("VAULT_TOKEN")
		Token:   os.Getenv("VAULT_ADDR")
	}, os.Getenv("MIGRATION_DIR"))
	if err != nil {
		log.Fatalf("Migration failed: %v", err)
	}
}

func migrate(vaultConfig vault.VaultConfig, migrateDir) error {
	// Load configuration from YAML file using Viper
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		return errors.Wrap(err, "failed to read config")
	}

	// Retrieve configuration values
	var conf AppConfig
	err = viper.Unmarshal(&conf)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshal config")
	}

	// Retrieve MySQL credentials from Vault
	// conf.Database.VaultMount
	creds, err := vault.GetSQLCredentials(vaultConfig, "/creds/sudo")
	if err != nil {
		return errors.Wrap(err, "failed to get credentials from Vault")
	}

	// Set MySQL credentials
	databaseURL := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
		creds.Username, creds.Password, conf.Database.Host, conf.Database.Port, conf.Database.Database)

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
