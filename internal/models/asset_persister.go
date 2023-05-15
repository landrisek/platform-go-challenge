package models

func (m *MySQLPersister) GetAssetWithExtendedType(assetID int64) (*Asset, error) {
	sql := `
		SELECT 
			asset.id,
			asset.type_id,
			asset.description,
			asset_type.type AS extended_type
		FROM asset
		INNER JOIN asset_type ON asset.type_id = asset_type.id
		WHERE asset.id = ?
	`
	var dbAsset Asset
	err := m.Conn.Get(&dbAsset, sql, assetID)
	if err != nil {
		return nil, err
	}
	asset := &Asset{
		ID:          dbAsset.ID,
		TypeID:      dbAsset.TypeID,
		Description: dbAsset.Description,
	}
	return asset, nil
}