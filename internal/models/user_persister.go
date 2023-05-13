package models

func (m *MySQLPersister) FindUserWithAssets(userID int64) (*models.User, error) {
	user, err := m.FindUser(userID)
	if err != nil {
		return nil, err
	}

	sql := `
		SELECT 
			asset.id,
			asset.type_id,
			asset.description,
			asset_type.type AS extended_type
		FROM asset
		INNER JOIN asset_type ON asset.type_id = asset_type.id
		INNER JOIN user_favorite ON user_favorite.asset_id = asset.id
		WHERE user_favorite.user_id = ?
	`
	var dbAssets []dbAssetWithExtendedType
	err = m.Conn.Select(&dbAssets, sql, userID)
	if err != nil {
		return nil, err
	}

	for _, dbAsset := range dbAssets {
		asset := &models.Asset{
			ID:          dbAsset.ID,
			TypeID:      dbAsset.TypeID,
			Description: dbAsset.Description,
			ExtendedType: dbAsset.ExtendedType,
		}
		user.Assets = append(user.Assets, asset)
	}

	return user, nil
}