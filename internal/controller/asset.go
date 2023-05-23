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
	"time"

	"github.com/landrisek/platform-go-challenge/internal/models"
	"github.com/landrisek/platform-go-challenge/internal/sagas"
	"github.com/landrisek/platform-go-challenge/internal/vault"

	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
)

func RunAsset(ctx context.Context, vaultConfig vault.VaultConfig, dbConfig models.DBConfig, redisAddr, blacklistAddr, port string) error {
	// redis
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// sql
	db, err := models.OpenDB(vaultConfig, dbConfig)
	if err != nil {
		return err
	}

	// sagas
	createOrchestrator := sagas.Create(db, blacklistAddr)
	readOrchestrator := sagas.Read(db, redisClient)
	updateOrchestrator := sagas.Update(db)
	deleteOrchestartor := sagas.Delete(db)

	router := mux.NewRouter()

	// Define the CRUD routes
	router.HandleFunc("/create", Generic(createOrchestrator.Run, db, "create")).Methods(http.MethodPost)
	router.HandleFunc("/read", Generic(readOrchestrator.Run, db, "read")).Methods(http.MethodPost)
	router.HandleFunc("/update", Generic(updateOrchestrator.Run, db, "update")).Methods(http.MethodPut)
	router.HandleFunc("/delete", Generic(deleteOrchestartor.Run, db, "delete")).Methods(http.MethodDelete)

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

func Generic(orchestratorFn orchestratorFunc, db *sqlx.DB, permission string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		err := Authenticate(request.Header, db, permission)
		if err != nil {
			http.Error(writer, "Unauthorized", http.StatusUnauthorized)
			return
		}
		// Read the request body
		body, err := ioutil.ReadAll(request.Body)
		if err != nil {
			log.Println("Error on reading request:", err)
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}

		// Create a new GenericRequest instance
		var genericReq sagas.GenericRequest

		// Unmarshal the JSON body into the GenericRequest struct
		err = json.Unmarshal(body, &genericReq.Data)
		if err != nil {
			log.Println("Error on unmarshaling request in generics: ", err)
			http.Error(writer, "Bad Request", http.StatusBadRequest)
			return
		}

		genericResp, err := orchestratorFn(genericReq)
		if err != nil {
			log.Println("Error on running orchestrator:", err)
			http.Error(writer, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Convert the response to JSON
		jsonResponse, err := json.Marshal(genericResp)
		if err != nil {
			log.Println("Error on marshalling response:", err)
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

// Log logs an error message along with the provided error.
func Log(err error, msg string) {
	if err != nil {
		log.Fatalf(msg+": %s", err)
	}
}
