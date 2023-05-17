// +build unit

package repository

import (
	"context"
	"errors"
	"sync"
	"testing"

	"github.com/go-redis/redis"
	"github.com/landrisek/platform-go-challenge/internal/repository"
)

func TestFindBlacklist(t *testing.T) {
	redisClient := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", 
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup

	errChan := make(chan error)

	data := map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"nested1": "nestedValue1",
			"nested2": "nestedValue2",
		},
	}

	// Add wait group for the initial call to findBlacklist
	wg.Add(1)

	// Call findBlacklist
	go repository.FindBlacklist(ctx, redisClient, data, &wg, errChan)

	// Wait for the completion of findBlacklist
	wg.Wait()

	// Check if any errors occurred
	select {
	case err := <-errChan:
		t.Fatalf("Unexpected error: %v", err)
	default:
		// No errors, continue with assertions
	}

	// Assert the result
	expectedData := map[string]interface{}{
		"key1": "value1",
		"key2": map[string]interface{}{
			"nested1": "nestedValue1",
			"nested2": "nestedValue2",
		},
	}

	if !isEqual(expectedData, data) {
		t.Errorf("Unexpected result. Expected: %+v, Got: %+v", expectedData, data)
	}
}

// Helper function to compare two maps recursively
func isEqual(expected, actual map[string]interface{}) bool {
	if len(expected) != len(actual) {
		return false
	}

	for key, expectedValue := range expected {
		actualValue, ok := actual[key]
		if !ok {
			return false
		}

		switch expectedValue := expectedValue.(type) {
		case string:
			actualValue, ok := actualValue.(string)
			if !ok || expectedValue != actualValue {
				return false
			}
		case map[string]interface{}:
			actualValue, ok := actualValue.(map[string]interface{})
			if !ok || !isEqual(expectedValue, actualValue) {
				return false
			}
		default:
			return false
		}
	}

	return true
}
