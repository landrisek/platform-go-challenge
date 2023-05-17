package utils

import (
	"testing"
	"time"
)

func TestExponentialBackoff(t *testing.T) {
	maxRetries := 5
	initialDelay := time.Millisecond * 100
	maxDelay := time.Second

	backoff := NewExponentialBackoff(maxRetries, initialDelay, maxDelay)

	for backoff.ShouldRetry() {
		if !backoff.Next() {
			t.Errorf("Expected Next to return true, got false")
		}
	}

	if backoff.ShouldRetry() {
		t.Errorf("Expected ShouldRetry to return false, got true")
	}
}

func TestExponentialBackoffWithMaxRetries(t *testing.T) {
	maxRetries := 3
	initialDelay := time.Millisecond * 100
	maxDelay := time.Second

	backoff := NewExponentialBackoff(maxRetries, initialDelay, maxDelay)

	for backoff.ShouldRetry() {
		if !backoff.Next() {
			t.Errorf("Expected Next to return true, got false")
		}
	}

	if backoff.ShouldRetry() {
		t.Errorf("Expected ShouldRetry to return false, got true")
	}
}

func TestExponentialBackoffWithMaxDelay(t *testing.T) {
	maxRetries := 5
	initialDelay := time.Millisecond * 100
	maxDelay := time.Millisecond * 500

	backoff := NewExponentialBackoff(maxRetries, initialDelay, maxDelay)

	for backoff.ShouldRetry() {
		if !backoff.Next() {
			t.Errorf("Expected Next to return true, got false")
		}
	}

	if backoff.ShouldRetry() {
		t.Errorf("Expected ShouldRetry to return false, got true")
	}
}
