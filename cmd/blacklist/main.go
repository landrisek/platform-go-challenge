package main

import (
	"fmt"
	"log"
	"os"

	"github.com/landrisek/platform-go-challenge/internal/controller"
)

func main() {
	redisAddr := fmt.Sprintf("%s:%s", os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"))
	errChan := make(chan error, 10) 
	defer close(errChan)
	// draining error channel
	go func() {
		for err := range errChan {
			log.Println("Received error:", err)
		}
	}()
	controller.Blacklist(redisAddr, os.Getenv("BLACKLIST_PORT"), errChan)
}
