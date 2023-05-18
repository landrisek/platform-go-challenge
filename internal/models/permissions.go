package models

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/go-sql-driver/mysql"
)

type Permission struct {
	ID     int    `db:"id"`
	Token  string `db:"token"`
	Create bool   `db:"create"`
	Read   bool   `db:"read"`
	Write  bool   `db:"write"`
	Delete bool   `db:"delete"`
}

func ReadPermission(db *sqlx.DB, id int) (*Permission, error) {
	permission := Permission{}
	err := db.Get(&permission, "SELECT * FROM permissions WHERE id = ?", id)
	if err != nil {
		return nil, err
	}
	return &permission, nil
}
