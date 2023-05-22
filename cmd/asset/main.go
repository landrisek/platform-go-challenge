package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/landrisek/platform-go-challenge/internal/controller"
	"github.com/landrisek/platform-go-challenge/internal/models"
	"github.com/landrisek/platform-go-challenge/internal/vault"
)

func main() {
	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	serverPort := os.Getenv("ASSET_PORT")
	vaultConfig := vault.VaultConfig{
		Address: os.Getenv("VAULT_ADDR"),
		Token:   os.Getenv("VAULT_TOKEN"),
		Mount: os.Getenv("VAULT_MOUNT"),
	}
	port, err := strconv.Atoi(os.Getenv("MYSQL_PORT"))
	if err != nil {
		log.Fatalf("Invalid port: %v", err)
	}
	dbConfig := models.DBConfig{
		Host:       os.Getenv("MYSQL_HOST"),
		Port:       port,
		Database:   os.Getenv("MYSQL_DATABASE"),
	}
	blacklistAddr := fmt.Sprintf("http://localhost:%s", os.Getenv("BLACKLIST_PORT"))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// HINT: in case number of input parameters will increase we will introduce specific structure for them
	err = controller.RunAsset(ctx, vaultConfig, dbConfig, redisAddr, blacklistAddr, serverPort)
	if err != nil {
		log.Println("Received error:", err)
	}
}