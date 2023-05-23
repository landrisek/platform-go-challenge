package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Chart struct {
	ID         int64  `db:"id"          json:"id"`
	Title      string `db:"title"       json:"title,omitempty"`
	AxesTitles string `db:"axes_titles" json:"axes_titles,omitempty"`
	Data       string `db:"data"        json:"data,omitempty"`
	// this will be not inserted in DB directly
	Description string `json:"description"`
	Error       string `json:"error"`
	AssetID     int64  `db:"asset_id"`
}

type ChartSafeUpdate struct {
	ID         int64   `db:"id"          json:"id"`
	Title      *string `db:"title"       json:"title,omitempty"`
	AxesTitles *string `db:"axes_titles" json:"axes_titles,omitempty"`
	Data       *string `db:"data"        json:"data,omitempty"`
	// this will be not inserted in DB directly
	Description *string `json:"description"`
	Error       string  `json:"error"`
	AssetID     int64   `db:"asset_id"`
}

const charts = "charts"

func CreateChart(db *sqlx.DB, chart Chart, userID int64) error {
	assetID, err := createAsset(db, Asset{
		Description: chart.Description,
		Type:        charts,
		UserID:      userID,
	})
	if err != nil {
		return err
	}

	chartInsertQuery := fmt.Sprintf("INSERT INTO %s (assets_id, title, axes_titles, data) VALUES (?, ?, ?, ?)", charts)
	_, err = db.Exec(chartInsertQuery, assetID, chart.Title, chart.AxesTitles, chart.Data)
	if err != nil {
		return err
	}
	return nil
}

func ReadCharts(db *sqlx.DB, userID int64) ([]Chart, error) {
	query := `
		SELECT ` + charts + `.id, ` + charts + `.axes_titles, ` + charts + `.data, ` + assets + `.description
		FROM ` + charts + ` 
		INNER JOIN ` + assets + ` ON ` + charts + `.assets_id = ` + assets + `.id
		WHERE ` + assets + `.user_id = ?
	`

	var items []Chart
	err := db.Select(&items, query, userID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

// UpdateChart is doing safe update.
// It update only those field values which are presented in given json
func UpdateChart(db *sqlx.DB, chart ChartSafeUpdate, userID int64) error {
	// Build the update query dynamically based on the non-nil fields in the chart struct
	var updateFields string
	var updateValues []interface{}

	if chart.Title != nil {
		updateFields += charts + ".title = ?, "
		updateValues = append(updateValues, *chart.Title)
	}

	if chart.AxesTitles != nil {
		updateFields += charts + ".axes_titles = ?, "
		updateValues = append(updateValues, *chart.AxesTitles)
	}

	if chart.Data != nil {
		updateFields += charts + ".data = ?, "
		updateValues = append(updateValues, *chart.Data)
	}

	if chart.Description != nil {
		updateFields += assets + ".description = ?, "
		// dereference
		updateValues = append(updateValues, *chart.Description)
	}

	if updateFields == "" {
		return fmt.Errorf("Chart with no valid attribute was provided")
	}

	// remove last comma
	updateFields = updateFields[:len(updateFields)-2]

	chartUpdateQuery := `
		UPDATE ` + charts + `
		INNER JOIN ` + assets + ` ON ` + charts + `.assets_id = ` + assets + `.id
		SET ` + updateFields + `
		WHERE ` + assets + `.user_id = ? AND ` + charts + `.id = ?
	`

	updateValues = append(updateValues, userID, chart.ID)

	result, err := db.Exec(chartUpdateQuery, updateValues...)
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

// DeleteChart delete both asset and char row due FK ON CASCADE delete
func DeleteChart(db *sqlx.DB, userID, chartID int64) error {
	assetQuery := `
		DELETE ` + assets + ` 
		FROM ` + assets + ` 
		INNER JOIN ` + charts + ` ON ` + charts + `.assets_id = ` + assets + `.id
		WHERE ` + assets + `.user_id = ? AND ` + charts + `.id = ?
	`

	result, err := db.Exec(assetQuery, userID, chartID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no rows deletedd on %s table", assets)
	}

	return nil
}
