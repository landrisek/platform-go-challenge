package models

type Asset struct {
	ID          int64  `db:"id"`
	Description string `db:"description"`
}