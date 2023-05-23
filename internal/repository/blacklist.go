package repository

import (
	"context"
	"sync"

	"github.com/go-redis/redis"
	"github.com/landrisek/platform-go-challenge/internal/models"
)

const blacklistPrefix = "blacklist"

func Blacklist(ctx context.Context, client *redis.Client, data []models.User, errChan chan<- error) []models.User {
	var wg sync.WaitGroup
	wg.Add(1)
	go findBlacklist(ctx, client, data, &wg, errChan)
	wg.Wait()
	return data
}

func findBlacklist(ctx context.Context, client *redis.Client, data []models.User, wg *sync.WaitGroup, errChan chan<- error) {
	defer wg.Done()

	select {
	case <-ctx.Done():
		return
	default:
	}

	// Set up the semaphore
	sem := make(chan struct{}, 10) // Change the value according to the desired concurrency limit

	for key := range data {
		user := &data[key]
		wg.Add(1)
		go func(user *models.User) {
			defer wg.Done()
			sem <- struct{}{} // Acquire a semaphore slot
			blacklistAssets(ctx, client, user, errChan)
			<-sem // Release the semaphore slot
		}()
	}
}

func blacklistAssets(ctx context.Context, client *redis.Client, user *models.User, errChan chan<- error) {

	var wg sync.WaitGroup
	wg.Add(len(user.Charts) + len(user.Insights) + len(user.Audiences))

	for i := range user.Audiences {
		go func(i int) {
			defer wg.Done()
			audience := &user.Audiences[i]
			if blacklisted, err := client.Get(blacklistPrefix + "." + audience.Description).Result(); err == nil {
				audience.Description = blacklisted
			} else {
				errChan <- err
			}
		}(i)
	}

	for i := range user.Charts {
		go func(i int) {
			defer wg.Done()
			chart := &user.Charts[i]
			if blacklisted, err := client.Get(blacklistPrefix + "." + chart.Description).Result(); err == nil {
				chart.Description = blacklisted
			} else {
				errChan <- err
			}
		}(i)
	}

	for i := range user.Insights {
		go func(i int) {
			defer wg.Done()
			insight := &user.Insights[i]
			if blacklisted, err := client.Get(blacklistPrefix + "." + insight.Description).Result(); err == nil {
				insight.Description = blacklisted
			} else {
				errChan <- err
			}
		}(i)
	}

	wg.Wait()
}
