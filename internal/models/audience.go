package models

import (
	"fmt"
	
	"github.com/jmoiron/sqlx"
)

type Audience struct {
	ID              int64  `db:"id" json:"id"`
	Characteristics string `db:"characteristics" json:"characteristics"`
	Description     string `db:"description" json:"description"`
}

const audiences = "audiences"

func CreateAudience(db *sqlx.DB, audience Audience, userID int64) error {
	assetID, err := createAsset(db, Asset{
		Description: audience.Description,
		Type: audiences,
		UserID: userID,
	})
	if err != nil {
		return err
	}

	audienceInsertQuery := fmt.Sprintf("INSERT INTO %s (assets_id, characteristics) VALUES (?, ?)", audiences)
	_, err = db.Exec(audienceInsertQuery, assetID, audience.Characteristics)
	if err != nil {
		return err
	}
	return nil
}