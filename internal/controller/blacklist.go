package controller

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/landrisek/platform-go-challenge/internal/repository"

	"github.com/go-redis/redis"
)

// HINT: channel is only for writing here, not draining
func Blacklist(redisAddr, blacklistPort string, errChan chan<- error) {
	// HINT: We can create this inside handleConnection, but we keep it outside, 
	// so we can share pointer to redis client between mulitple goroutines
	// if we decide to scale this way
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	// Create a context with cancellation
	ctx, cancel := context.WithCancel(context.Background())

	// Handle OS signals to trigger cancellation
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChan
		cancel()
	}()

	if runtime.NumCPU() > 1 {
		// HINT: if there are multiple CPUs, let`s be not too greedy and use less than the total available
		runtime.GOMAXPROCS(runtime.NumCPU() / 2)
	} else {
		runtime.GOMAXPROCS(1)
	}

	// Create a TCP server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", blacklistPort))
	if err != nil {
		errChan <- fmt.Errorf("failed to start server: %s", err)
		return
	}
	defer listener.Close()

	fmt.Println(fmt.Sprintf("Server started. Listening on port %s...", blacklistPort))

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			errChan <- fmt.Errorf("error accepting connection: %s", err)
			continue
		}

		// Handle each connection in a separate goroutine
		go handleConnection(ctx, conn, redisClient, errChan)
	}
}

func handleConnection(ctx context.Context, conn net.Conn, client *redis.Client, errChan chan<- error) {
	defer conn.Close()

	// Read client request
	reader := bufio.NewReader(conn)
	request, err := reader.ReadString('\n')
	if err != nil {
		errChan <- fmt.Errorf("error reading request: %s", err)
		return
	}

	// Parse the JSON request
	var jsonData map[string]interface{}
	err = json.Unmarshal([]byte(request), &jsonData)
	if err != nil {
		errChan <- fmt.Errorf("error parsing JSON: %s", err)
		return
	}

	// Process the JSON data
	repository.Blacklist(ctx, client, jsonData, errChan)
	
	// Convert the response back to JSON
	responseJSON, err := json.Marshal(jsonData)
	if err != nil {
		errChan <- fmt.Errorf("error encoding JSON response: %s", err)
		return
	}

	// Send the JSON response back to the client
	_, err = conn.Write(responseJSON)
	if err != nil {
		errChan <- fmt.Errorf("error sending response to client: %s", err)
		return
	}
}