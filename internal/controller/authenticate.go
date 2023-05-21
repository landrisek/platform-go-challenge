package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/landrisek/platform-go-challenge/internal/models"

	"github.com/jmoiron/sqlx"
)

func Authenticate(header http.Header, db *sqlx.DB, permission string) error {
	authHeader := header.Get("Authorization")
	if authHeader == "" {
		return fmt.Errorf("Unauthorized")
	}
	authParts := strings.Split(authHeader, " ")
	var authToken string
	if len(authParts) == 2 && authParts[0] == "Bearer" {
		authToken = authParts[1]
	}
	if authToken == "" {
		return fmt.Errorf("Empty token")
	}
	_, err := models.IsValidPermission(db, permission, authToken)
	return err
}
