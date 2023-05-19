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

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type CRUD struct {
	db    *sqlx.DB
	redis *redis.Client
}

func RunUser(vaultConfig vault.VaultConfig, dbConfig models.DBConfig, port, redisAddr string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// sql
	db, err := models.OpenDB(vaultConfig, dbConfig)
	if err != nil {
		return err
	}

	// user handler
	crud := CRUD{
		db:    db,
		redis: redisClient,
	}

	router := mux.NewRouter()

	// Define the CRUD routes
	router.HandleFunc("/create", crud.Create).Methods(http.MethodPost)
	// TODO: implement rest of RUD operations

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
	permissionID, err := Authenticate(request.Header, crud.redis)
	if err != nil {
		http.Error(writer, "Unauthorized", http.StatusUnauthorized)
		return	
	}
	permission, err := models.ReadPermission(crud.db, permissionID)
	if err != nil {
		http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	if !permission.Create {
		http.Error(writer, "Unauthorized: Insufficient permission", http.StatusForbidden)
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

	// Encode the response as JSON
	responseJSON, err := json.Marshal(struct {
		UserID int64 `json:"user_id"`
	}{
		UserID: userID,
	})
	if err != nil {
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusOK)
	writer.Write(responseJSON)

}
