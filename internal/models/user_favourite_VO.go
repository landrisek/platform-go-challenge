package models

type UserFavorite struct {
	UserID  int64 `db:"user_id"`
	AssetID int64 `db:"asset_id"`
}
