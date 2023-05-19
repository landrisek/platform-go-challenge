package models

import (
	"fmt"
	
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Charts    []Chart    `json:"charts"`
	Insights  []Insight  `json:"insights"`
	Audiences []Audience `json:"audiences"`
}

func CreateUser(db *sqlx.DB, user User) (int64, error) {
	// Insert user data into the users table
	userInsertQuery := fmt.Sprintf("INSERT INTO users (id, name) VALUES (%d, '%s')", user.ID, user.Name)
	result, err := db.Exec(userInsertQuery)
	if err != nil {
		return 0, err
	}

	insertedID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return insertedID, nil
}