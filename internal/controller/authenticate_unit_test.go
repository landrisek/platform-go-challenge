//go:build unit
// +build unit

package controller

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func TestAuthenticate(t *testing.T) {
	db, err := sqlx.Connect("sqlite3", "test.db")
	if err != nil {
		t.Fatal("Failed to connect to the database:", err)
	}

	testCases := []struct {
		name        string
		header      http.Header
		permission  string
		expectedErr error
	}{
		{
			name:        "EmptyHeader",
			header:      make(http.Header),
			permission:  "read",
			expectedErr: fmt.Errorf("Unauthorized"),
		},
		{
			name: "EmptyToken",
			header: http.Header{
				"Authorization": []string{"Bearer"},
			},
			permission:  "read",
			expectedErr: fmt.Errorf("Empty token"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := Authenticate(tc.header, db, tc.permission)
			if err == nil {
				t.Error("Expected an error, but got nil")
			} else if err.Error() != tc.expectedErr.Error() {
				t.Errorf("Unexpected error. Expected: %v, Got: %v", tc.expectedErr, err)
			}
		})
	}
}
