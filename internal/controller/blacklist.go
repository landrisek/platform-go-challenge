package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/landrisek/platform-go-challenge/internal/models"
	"github.com/landrisek/platform-go-challenge/internal/repository"

	"github.com/go-redis/redis"
)

func Blacklist(redisAddr, blacklistPort string, errChan chan<- error) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	ctx, cancel := context.WithCancel(context.Background())

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()

	if runtime.NumCPU() > 1 {
		runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	} else {
		runtime.GOMAXPROCS(1)
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		requestBody, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errChan <- fmt.Errorf("error reading request body: %s", err)
			return
		}

		// HINT: this will enrich return generic response with all field in given types
		// even if they are not provided
		var jsonData []models.User
		err = json.Unmarshal(requestBody, &jsonData)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			errChan <- fmt.Errorf("error parsing JSON: %s", err)
			return
		}

		jsonData = repository.Blacklist(ctx, redisClient, jsonData, errChan)

		responseJSON, err := json.Marshal(jsonData)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			errChan <- fmt.Errorf("error encoding JSON response: %s", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, err = w.Write(responseJSON)
		if err != nil {
			errChan <- fmt.Errorf("error sending response: %s", err)
			return
		}
	})

	err := http.ListenAndServe(fmt.Sprintf(":%s", blacklistPort), nil)
	if err != nil {
		errChan <- fmt.Errorf("failed to start server: %s", err)
	}
}
