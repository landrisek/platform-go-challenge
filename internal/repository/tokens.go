package repository

import (
	"github.com/go-redis/redis"
)

const tokensPrefix = "tokens"

func IsValidToken(client *redis.Client, token string) (string, error) {
	return client.Get(tokensPrefix + "." + token).Result()
}
