package repository

import (
	"context"
	"sync"

	"github.com/go-redis/redis"
)

const blacklistPrefix = "blacklist"

func Blacklist(ctx context.Context, client *redis.Client, data map[string]interface{}, errChan chan<- error) {
	var wg sync.WaitGroup
	wg.Add(1)
	go findBlacklist(ctx, client, data, &wg, errChan)
	wg.Wait()
}

func findBlacklist(ctx context.Context, client *redis.Client, data map[string]interface{}, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default:
	}

	for key, value := range data {
		if str, ok := value.(string); ok {
			blacklisted, err := client.Get(blacklistPrefix + "." + str).Result()
			if err == nil {
				data[key] = blacklisted
			} else {
				errChan <- err
			}
		} else if nestedData, ok := value.(map[string]interface{}); ok {
			wg.Add(1)
			// TODO: goroutine pool managment
			go findBlacklist(ctx, client, nestedData, wg, errChan)
		}
	}
}
