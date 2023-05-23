package models

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type User struct {
	ID        int64      `json:"id"`
	Name      string     `json:"name"`
	Audiences []Audience `json:"audiences"`
	Charts    []Chart    `json:"charts"`
	Insights  []Insight  `json:"insights"`
}

type UserSafeUpdate struct {
	ID        int64                `json:"id"`
	Name      string               `json:"name"`
	Audiences []AudienceSafeUpdate `json:"audiences"`
	Charts    []ChartSafeUpdate    `json:"charts"`
	Insights  []InsightSafeUpdate  `json:"insights"`
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
