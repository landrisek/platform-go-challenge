package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/landrisek/platform-go-challenge/internal/repository"

	"github.com/go-redis/redis"
)

func Authenticate(header http.Header, client *redis.Client) (int, error) {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return 0, fmt.Errorf("Unauthorized")
	}
	authParts := strings.Split(authHeader, " ")
	var authToken string
	if len(authParts) == 2 && authParts[0] == "Bearer" {
		authToken = authParts[1]
	}
	if authToken == "" {
		return 0, fmt.Errorf("Empty token")
	}
	tokenID, err := repository.IsValidToken(client, authToken)
	if err != nil  {
		return 0, err
	}
	id, err := strconv.Atoi(tokenID)
	return id, nil
}
