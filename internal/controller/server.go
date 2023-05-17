package controller

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/landrisek/platform-go-challenge/internal/models"
	"github.com/landrisek/platform-go-challenge/internal/repository"
	"github.com/landrisek/platform-go-challenge/internal/vault"

	"github.com/pkg/errors"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type CRUD struct {
	db *sqlx.DB
}

func NewCRUD(vaultConfig vault.VaultConfig, dbConfig models.DBConfig) (*CRUD, error) {
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

	return &CRUD{
		db: db,
	}, nil
}

func RunServer(vaultConfig vault.VaultConfig, dbConfig models.DBConfig, port, redisAddr string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	crud, err := NewCRUD(vaultConfig, dbConfig)
	if err != nil {
		return err
	}

	router := mux.NewRouter()

	// Define the CRUD routes
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})
	
	router.HandleFunc("/read/{user_id}", Authenticate(crud.Read, redisClient)).Methods(http.MethodGet)
	router.HandleFunc("/create", Authenticate(crud.Create, redisClient)).Methods(http.MethodPost)
	router.HandleFunc("/update", Authenticate(crud.Update, redisClient)).Methods(http.MethodPut)
	router.HandleFunc("/delete", Authenticate(crud.Delete, redisClient)).Methods(http.MethodDelete)

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

func Authenticate(route http.HandlerFunc, client *redis.Client) http.HandlerFunc {
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
		if request.Method == http.MethodGet {
			// TODO: use cache
		}
		// run route if authentication is successful
		route(writer, request)
	}
}

func Cache(route http.HandlerFunc, client *redis.Client) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		// todo
	}
}


// HINT: explain why one pointer, one not
func (crud *CRUD) Read(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("-------read-------")
}

func (crud *CRUD) Create(writer http.ResponseWriter, request *http.Request) {

}

func (crud *CRUD) Update(writer http.ResponseWriter, request *http.Request) {

}

func (crud *CRUD) Delete(writer http.ResponseWriter, request *http.Request) {

}

// Log logs an error message along with the provided error.
func Log(err error, msg string) {
	if err != nil {
		log.Fatalf(msg+": %s", err)
	}
}
