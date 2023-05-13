package models

type Asset struct {
	ID          int64  `db:"id"`
	TypeID      int64  `db:"type_id"`
	Description string `db:"description"`
}
