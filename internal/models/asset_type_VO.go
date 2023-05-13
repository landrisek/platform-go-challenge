package models

type AssetType struct {
	ID   int64  `db:"id"`
	Type string `db:"type"`
}