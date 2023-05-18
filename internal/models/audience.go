package models

type Audience struct {
	ID             int64  `db:"id" json:"id"`
	Characteristics string `db:"characteristics" json:"characteristics"`
	Description    string `db:"description" json:"description"`
}