package models

import (
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Permission struct {
	ID     int    `db:"id"`
	Token  string `db:"token"`
	Create bool   `db:"create"`
	Read   bool   `db:"read"`
	Update bool   `db:"update"`
	Delete bool   `db:"delete"`
}

func IsValidPermission(db *sqlx.DB, column string, token string) (*Permission, error) {
	permission := Permission{}
	err := db.Get(&permission, fmt.Sprintf("SELECT * FROM permissions WHERE `%s` = ? AND token = ?", column), 1, token)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}
