package models

type Insight struct {
	ID          int64  `db:"id" json:"id"`
	Text        string `db:"text" json:"text"`
	Description string `db:"description" json:"description"`
}