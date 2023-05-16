package models

import (
	"github.com/jmoiron/sqlx"
)

type Chart struct {
	ID         int64  `db:"id"`
	Title      string `db:"title"`
	AxesTitles string `db:"axes_titles"`
	Data       string `db:"data"`
}

const (
	assets = "assets"
	charts = "charts"
)

func GetChartsByUserID(db *sqlx.DB, userID int) ([]Chart, error) {
	query := `
		SELECT ` + charts + `.id, ` + charts + `.axes_titles, ` + charts + `.data, ` + charts + `.description
		FROM ` + charts + ` 
		INNER JOIN ` + assets + ` ON ` + charts + `.id = ` + assets + `.id
		WHERE a.user_id = ?
	`

	var charts []Chart
	err := db.Select(&charts, query, userID)
	if err != nil {
		return nil, err
	}

	return charts, nil
}