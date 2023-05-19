package models

import (
	"fmt"
	
	"github.com/jmoiron/sqlx"
)

type Chart struct {
	ID          int64  `db:"id" json:"id"`
	Description string `db:"description" json:"description"`
	Title       string `db:"title" json:"title"`
	AxesTitles  string `db:"axes_titles" json:"axes_titles"`
	Data        string `db:"data" json:"data"`
	// FK
	AssetID     int64  `db:"asset_id"`
}

const charts = "charts"

func CreateChart(db *sqlx.DB, chart Chart, userID int64) error {
	assetID, err := createAsset(db, Asset{
		Description: chart.Description,
		Type: charts,
		UserID: userID,
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

func ReadCharts(db *sqlx.DB, userID int) ([]Chart, error) {
	query := `
		SELECT ` + charts + `.id, ` + charts + `.axes_titles, ` + charts + `.data, ` + charts + `.description
		FROM ` + charts + ` 
		INNER JOIN ` + assets + ` ON ` + charts + `.id = ` + assets + `.id
		WHERE ` + assets + `.user_id = ?
	`

	var charts []Chart
	err := db.Select(&charts, query, userID)
	if err != nil {
		return nil, err
	}

	return charts, nil
}


