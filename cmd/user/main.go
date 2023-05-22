package main

import (
	"context"
	"log"
	"os"
	"strconv"

	"github.com/landrisek/platform-go-challenge/internal/controller"
	"github.com/landrisek/platform-go-challenge/internal/models"
	"github.com/landrisek/platform-go-challenge/internal/vault"
)

func main() {
	serverPort := os.Getenv("USER_PORT")
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

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	err = controller.RunUser(ctx, vaultConfig, dbConfig, serverPort)
	if err != nil {
		log.Println("Received error:", err)
	}
}