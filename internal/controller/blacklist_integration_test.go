// +build integration

package controller

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/go-redis/redis"
	"github.com/landrisek/platform-go-challenge/internal/controller"
)

// MockRedisClient is a mock implementation of the Redis client for testing.
type MockRedisClient struct {
	GetFunc func(key string) *redis.StringCmd
}

// Get is a mock implementation of the Redis client's Get method.
func (m *MockRedisClient) Get(key string) *redis.StringCmd {
	if m.GetFunc != nil {
		return m.GetFunc(key)
	}
	return nil
}

// MockNetConn is a mock implementation of the net.Conn interface for testing.
type MockNetConn struct {
	ReadFunc  func([]byte) (int, error)
	WriteFunc func([]byte) (int, error)
	CloseFunc func() error
}

// Read is a mock implementation of the Read method.
func (m *MockNetConn) Read(b []byte) (int, error) {
	if m.ReadFunc != nil {
		return m.ReadFunc(b)
	}
	return 0, nil
}

// Write is a mock implementation of the Write method.
func (m *MockNetConn) Write(b []byte) (int, error) {
	if m.WriteFunc != nil {
		return m.WriteFunc(b)
	}
	return 0, nil
}

// Close is a mock implementation of the Close method.
func (m *MockNetConn) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func TestBlacklist(t *testing.T) {
	// Create a Redis test server
	redisServer := mockredis.Start()
	defer redisServer.Close()

	// Create a Redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr: redisServer.Addr(),
	})

	// Start the blacklist server
	blacklistPort := "1234"
	errChan := make(chan error)
	go Blacklist(redisServer.Addr(), blacklistPort, errChan)

	// Wait for the server to start
	time.Sleep(100 * time.Millisecond)

	// Test cases
	tests := []struct {
		name           string
		request        string
		expectedResult string
		expectedError  error
	}{
		{
			name:           "ValidRequest",
			request:        `{"key": "value"}`,
			expectedResult: `{"key": "value"}`,
			expectedError:  nil,
		},
		{
			name:           "InvalidRequest",
			request:        `{"invalid}`,
			expectedResult: "",
			expectedError:  fmt.Errorf("error parsing JSON"),
		},
	}

	// Run test cases
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a TCP connection
			conn, err := net.Dial("tcp", "localhost:"+blacklistPort)
			if err != nil {
				t.Fatalf("Failed to create TCP connection: %v", err)
			}
			defer conn.Close()

			// Send request to the server
			_, err = fmt.Fprintf(conn, tc.request+"\n")
			if err != nil {
				t.Fatalf("Failed to send request to server: %v", err)
			}

			// Read response from the server
			response, err := bufio.NewReader(conn).ReadString('\n')
			if err != nil {
				t.Fatalf("Failed to read response from server: %v", err)
			}

			// Remove trailing newline character
			response = strings.TrimSuffix(response, "\n")

			// Assert the response
			if response != tc.expectedResult {
				t.Errorf("Unexpected response. Expected: %s, Got: %s", tc.expectedResult, response)
			}

			// Assert the error
			select {
			case err := <-errChan:
				if tc.expectedError == nil {
					t.Errorf("Unexpected error: %v", err)
				} else if err.Error() != tc.expectedError.Error() {
					t.Errorf("Unexpected error. Expected: %v, Got: %v", tc.expectedError, err)
				}
			default:
				if tc.expectedError != nil {
					t.Errorf("Expected error: %v, but no error occurred", tc.expectedError)
				}
			}
		})
	}

	// Wait for the server to shut down
	time.Sleep(100 * time.Millisecond)
}
