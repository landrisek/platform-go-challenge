package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/landrisek/platform-go-challenge/internal/models"
	"github.com/landrisek/platform-go-challenge/internal/vault"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

type CRUD struct {
	db *sqlx.DB
}

func RunUser(ctx context.Context, vaultConfig vault.VaultConfig, dbConfig models.DBConfig, port string) error {

	// sql
	db, err := models.OpenDB(vaultConfig, dbConfig)
	if err != nil {
		return err
	}

	// user handler
	crud := CRUD{
		db: db,
	}

	router := mux.NewRouter()

	// Define the CRUD routes
	router.HandleFunc("/create", crud.Create).Methods(http.MethodPost)
	// TODO: implement rest of (C)RUD operations

	server := &http.Server{
		Addr:           fmt.Sprintf(":%s", port),
		ReadTimeout:    60 * time.Second,
		WriteTimeout:   60 * time.Second,
		MaxHeaderBytes: 1 << 20,
		Handler:        router,
	}

	go func() {
		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			Log(err, "Error running tag server")
		}
	}()

	// HINT: Wait for SIGINT (Ctrl+C) signal
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint

	log.Println("Shutting down server...")
	if err := server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (crud CRUD) Create(writer http.ResponseWriter, request *http.Request) {
	err := Authenticate(request.Header, crud.db, "create")
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return
	}
	var user models.User
	err = json.NewDecoder(request.Body).Decode(&user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusBadRequest)
		return
	}
	userID, err := models.CreateUser(crud.db, user)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	user.ID = userID
	// Encode the response as JSON
	responseJSON, err := json.Marshal(user)
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseJSON)

}
