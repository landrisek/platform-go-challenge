package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Asset struct {
	ID          int64  `db:"id"`
	Description string `db:"description"`
	Type        string `db:"type"`
	// FK
	UserID int64 `db:"user_id"`
}

type AssetDescription string

const assets = "assets"

func createAsset(db *sqlx.DB, asset Asset) (int64, error) {
	assetInsertQuery := fmt.Sprintf("INSERT INTO %s (description, type, user_id) VALUES (?, ?, ?)", assets)
	res, err := db.Exec(assetInsertQuery, asset.Description, asset.Type, asset.UserID)
	if err != nil {
		return 0, err
	}
	assetID, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return assetID, nil
}
