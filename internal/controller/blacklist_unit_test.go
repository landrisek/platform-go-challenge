// +build unit

package controller

import (
	"context"
	"net"
	"testing"
)

func TestHandleConnection(t *testing.T) {
	// Create a mock net.Conn
	mockConn := &MockNetConn{
		ReadFunc: func(b []byte) (int, error) {
			// Simulate reading the client request
			copy(b, []byte(`{"key": "value"}`))
			return len(`{"key": "value"}`), nil
		},
		WriteFunc: func(b []byte) (int, error) {
			// Simulate writing the response back to the client
			expectedResponse := []byte(`{"key": "value"}`)
			if string(b) != string(expectedResponse) {
				t.Errorf("Unexpected response. Expected: %s, Got: %s", expectedResponse, b)
			}
			return len(b), nil
		},
		CloseFunc: func() error {
			return nil
		},
	}

	// Create a mock Redis client
	mockRedisClient := &MockRedisClient{}

	// Create a context
	ctx := context.Background()

	// Call the handleConnection function
	handleConnection(ctx, mockConn, mockRedisClient, nil)
	// Perform assertions on the mockConn
	if mockConn.ReadCalled != true {
		t.Errorf("Read method not called on the connection")
	}

	if mockConn.WriteCalled != true {
		t.Errorf("Write method not called on the connection")
	}

	if mockConn.CloseCalled != true {
		t.Errorf("Close method not called on the connection")
	}
}

