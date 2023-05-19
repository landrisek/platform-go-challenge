package models

import (
	"fmt"
	
	"github.com/jmoiron/sqlx"
)

type Insight struct {
	ID          int64  `db:"id" json:"id"`
	Text        string `db:"text" json:"text"`
	Description string `db:"description" json:"description"`
}

const insights = "insights"

func CreateInsight(db *sqlx.DB, insight Insight, userID int64) error {
	assetID, err := createAsset(db, Asset{
		Description: insight.Description,
		Type: insights,
		UserID: userID,
	})
	if err != nil {
		return err
	}
	insightInsertQuery := fmt.Sprintf("INSERT INTO %s (assets_id, text) VALUES (?, ?)", insights)
	_, err = db.Exec(insightInsertQuery, assetID, insight.Text)
	if err != nil {
		return err
	}
	return nil
}