package utils

import (
	"time"
)

type ExponentialBackoff struct {
	maxRetries     int
	initialDelay   time.Duration
	maxDelay       time.Duration
	currentRetries int
	currentDelay   time.Duration
}

func NewExponentialBackoff(maxRetries int, initialDelay, maxDelay time.Duration) *ExponentialBackoff {
	return &ExponentialBackoff{
		maxRetries:     maxRetries,
		initialDelay:   initialDelay,
		maxDelay:       maxDelay,
		currentRetries: 0,
		currentDelay:   initialDelay,
	}
}

func (b *ExponentialBackoff) Next() bool {
	if b.currentRetries >= b.maxRetries {
		return false
	}

	time.Sleep(b.currentDelay)
	b.currentRetries++
	b.currentDelay = b.currentDelay * 2

	if b.currentDelay > b.maxDelay {
		b.currentDelay = b.maxDelay
	}

	return true
}

func (b *ExponentialBackoff) ShouldRetry() bool {
	return b.currentRetries < b.maxRetries
}
