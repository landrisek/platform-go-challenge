package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Insight struct {
	ID   int64  `db:"id"          json:"id"`
	Text string `db:"text"        json:"text"`
	// this will be not inserted in DB directly
	Description string `json:"description"`
	Error       string `json:"error"`
}

type InsightSafeUpdate struct {
	ID   int64   `db:"id"          json:"id"`
	Text *string `db:"text"        json:"text"`
	// this will be not inserted in DB directly
	Description *string `json:"description"`
	Error       string  `json:"error"`
}

const insights = "insights"

// CreateInsight check if at least one valid provided json field is valid
// and then process it, otherwise it will error out
func CreateInsight(db *sqlx.DB, insight Insight, userID int64) error {
	assetID, err := createAsset(db, Asset{
		Description: insight.Description,
		Type:        insights,
		UserID:      userID,
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

func ReadInsights(db *sqlx.DB, userID int64) ([]Insight, error) {
	query := `
		SELECT ` + insights + `.id, ` + insights + `.text, ` + assets + `.description
		FROM ` + insights + ` 
		INNER JOIN ` + assets + ` ON ` + insights + `.assets_id = ` + assets + `.id
		WHERE ` + assets + `.user_id = ?
	`

	var items []Insight
	err := db.Select(&items, query, userID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// UpdateInsight is doing safe update.
// It update only those field values which are presented in given json
func UpdateInsight(db *sqlx.DB, insight InsightSafeUpdate, userID int64) error {
	// Build the update query dynamically based on the non-nil fields in the insight struct
	var updateFields string
	var updateValues []interface{}

	if insight.Text != nil {
		updateFields += insights + ".text = ?, "
		updateValues = append(updateValues, *insight.Text)
	}

	if insight.Description != nil {
		updateFields += assets + ".description = ?, "
		updateValues = append(updateValues, *insight.Description)
	}

	if updateFields == "" {
		return fmt.Errorf("Insight with no valid attribute was provided")
	}
	// remove last comma
	updateFields = updateFields[:len(updateFields)-2]

	insightUpdateQuery := `
		UPDATE ` + insights + `
		INNER JOIN ` + assets + ` ON ` + insights + `.assets_id = ` + assets + `.id
		SET ` + updateFields + `
		WHERE ` + assets + `.user_id = ? AND ` + insights + `.id = ?
	`

	updateValues = append(updateValues, userID, insight.ID)

	result, err := db.Exec(insightUpdateQuery, updateValues...)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows updated")
	}

	return nil
}

// DeleteInsight delete both asset and char row due FK ON CASCADE delete
func DeleteInsight(db *sqlx.DB, userID, insightID int64) error {
	assetQuery := `
		DELETE ` + assets + ` 
		FROM ` + assets + ` 
		INNER JOIN ` + insights + ` ON ` + insights + `.assets_id = ` + assets + `.id
		WHERE ` + assets + `.user_id = ? AND ` + insights + `.id = ?
	`

	result, err := db.Exec(assetQuery, userID, insightID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows deleted on %s table", assets)
	}

	return nil
}
