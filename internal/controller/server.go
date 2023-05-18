package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/landrisek/platform-go-challenge/internal/models"
	"github.com/landrisek/platform-go-challenge/internal/repository"
	"github.com/landrisek/platform-go-challenge/internal/sagas"
	"github.com/landrisek/platform-go-challenge/internal/vault"

	"github.com/pkg/errors"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

func openDB(vaultConfig vault.VaultConfig, dbConfig models.DBConfig) (*sqlx.DB, error) {
	// Define the MySQL database connection string

	dbURL, err := models.GetDatabaseURL(vaultConfig, dbConfig)
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

func RunServer(vaultConfig vault.VaultConfig, dbConfig models.DBConfig, port, redisAddr string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// sql
	db, err := openDB(vaultConfig, dbConfig)
	if err != nil {
		return err
	}
	// sagas
	createSagas := sagas.Create(db)

	router := mux.NewRouter()

	// Define the CRUD routes
	router.HandleFunc("/create", Authenticate(createSagas.Run, redisClient)).Methods(http.MethodPost)
	//router.HandleFunc("/read/{user_id}", Authenticate(create, redisClient)).Methods(http.MethodGet)
	//router.HandleFunc("/update", Authenticate(Update, redisClient)).Methods(http.MethodPut)
	//router.HandleFunc("/delete", Authenticate(Delete, redisClient)).Methods(http.MethodDelete)

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

type orchestratorFunc func(sagas.GenericRequest) (sagas.GenericResponse, error)

func Authenticate(orchestratorFn orchestratorFunc, client *redis.Client) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		authHeader := request.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		authParts := strings.Split(authHeader, " ")
		var authToken string
		if len(authParts) == 2 && authParts[0] == "Bearer" {
			authToken = authParts[1]
		}

		if authToken == "" || !repository.IsValidToken(client, authToken)  {
		    http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Read the request body
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}
		
		// Create a new GenericRequest instance
		var genericReq sagas.GenericRequest
		
		// Unmarshal the JSON body into the GenericRequest struct
		err = json.Unmarshal(body, &genericReq)
		if err != nil {
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}

		genericResp, err := orchestratorFn(genericReq)
		if err != nil {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Convert the response to JSON
		jsonResponse, err := json.Marshal(genericResp)
		if err != nil {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Set the Content-Type header to application/json
		writer.Header().Set("Content-Type", "application/json")

		// Set the response status code to 200
		writer.WriteHeader(http.StatusOK)

		// Write the JSON response to the writer
		_, err = writer.Write(jsonResponse)
		if err != nil {
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		
	}
}

// HINT: explain why one pointer, one not
func read(writer http.ResponseWriter, request *http.Request) {
	// TODO: cache
	fmt.Println("-------read-------")
}

func create(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("-------create-------")
}

func Update(writer http.ResponseWriter, request *http.Request) {

}

func Delete(writer http.ResponseWriter, request *http.Request) {

}

// Log logs an error message along with the provided error.
func Log(err error, msg string) {
	if err != nil {
		log.Fatalf(msg+": %s", err)
	}
}
