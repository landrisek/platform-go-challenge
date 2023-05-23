package models

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Audience struct {
	ID              int64  `db:"id"               json:"id"`
	Characteristics string `db:"characteristics" json:"characteristics,omitempty"`
	// this will be not inserted in DB directly
	Description string `json:"description,omitempty"`
	Error       string `json:"error"`
}

// structure with pointer was created to use heap memory
// where it is usefull - to omit not provided values
type AudienceSafeUpdate struct {
	ID              int64   `db:"id"               json:"id"`
	Characteristics *string `db:"characteristics" json:"characteristics,omitempty"`
	// this will be not inserted in DB directly
	Description *string `json:"description,omitempty"`
	Error       string  `json:"error"`
}

const audiences = "audiences"

func CreateAudience(db *sqlx.DB, audience Audience, userID int64) error {
	assetID, err := createAsset(db, Asset{
		Description: audience.Description,
		Type:        audiences,
		UserID:      userID,
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

func ReadAudiences(db *sqlx.DB, userID int64) ([]Audience, error) {
	query := `
		SELECT ` + audiences + `.id, ` + audiences + `.characteristics, ` + assets + `.description
		FROM ` + audiences + ` 
		INNER JOIN ` + assets + ` ON ` + audiences + `.assets_id = ` + assets + `.id
		WHERE ` + assets + `.user_id = ?
	`

	var items []Audience
	err := db.Select(&items, query, userID)
	if err != nil {
		return nil, err
	}

	return items, nil
}

func ReadAudience(db *sqlx.DB, userID int64, audienceID int64) (Audience, error) {
	query := `
		SELECT ` + audiences + `.id, ` + audiences + `.characteristics, ` + assets + `.description
		FROM ` + audiences + ` 
		INNER JOIN ` + assets + ` ON ` + audiences + `.assets_id = ` + assets + `.id
		WHERE ` + assets + `.user_id = ? AND ` + audiences + `.id = ?
	`

	var item Audience
	err := db.Get(&item, query, userID, audienceID)
	if err != nil {
		return Audience{}, err
	}

	return item, nil
}

// UpdateAudience is doing safe update.
// It update only those field values which are presented in given json
func UpdateAudience(db *sqlx.DB, audience AudienceSafeUpdate, userID int64) error {
	// Build the update query dynamically based on the non-nil fields in the audience struct
	var updateFields string
	var updateValues []interface{}

	if audience.Characteristics != nil {
		updateFields += audiences + ".characteristics = ?, "
		updateValues = append(updateValues, *audience.Characteristics)
	}

	if audience.Description != nil {
		updateFields += assets + ".description = ?, "
		updateValues = append(updateValues, *audience.Description)
	}

	if updateFields == "" {
		return fmt.Errorf("Audience with no valid attribute was provided")
	}

	// remove last comma
	updateFields = updateFields[:len(updateFields)-2]

	audienceUpdateQuery := `
		UPDATE ` + audiences + `
		INNER JOIN ` + assets + ` ON ` + audiences + `.assets_id = ` + assets + `.id
		SET ` + updateFields + `
		WHERE ` + assets + `.user_id = ? AND ` + audiences + `.id = ?
	`

	updateValues = append(updateValues, userID, audience.ID)

	result, err := db.Exec(audienceUpdateQuery, updateValues...)
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

// DeleteAudience delete both asset and char row due FK ON CASCADE delete
func DeleteAudience(db *sqlx.DB, userID, audienceID int64) error {
	assetQuery := `
		DELETE ` + assets + ` 
		FROM ` + assets + ` 
		INNER JOIN ` + audiences + ` ON ` + audiences + `.assets_id = ` + assets + `.id
		WHERE ` + assets + `.user_id = ? AND ` + audiences + `.id = ?
	`

	result, err := db.Exec(assetQuery, userID, audienceID)
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
