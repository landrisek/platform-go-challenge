package repository

import (
	"github.com/go-redis/redis"
)

const tokensPrefix = "tokens"

func IsValidToken(client *redis.Client, token string) bool {
	_, err := client.Get(tokensPrefix + "." + token).Result()
	if err == nil {
		return true
	}
	return false
}
