package models

import (
	"time"

	"github.com/pkg/errors"
)

// HINT: we can use this as value object, note implemented methods
// TODO: guid?
type User struct {
	ID     int64    `db:"id"`
	Name   string   `db:"name"`
	Email  string   `db:"email"`
	Assets []*Asset `db:"-"`
}

func (u *User) AddFavoriteAsset(assetID int64) {

}

// TODO: do encryption
func (u *User) AddEmail(email string) {
	u.Email = email
}

func (u *User) RemoveFavoriteAsset(assetID int64) error {
	return errors.New("Asset not found in favorites")
}

// GetFavoriteAssets returns a slice of asset IDs
func (u *User) GetFavoriteAssets() []int64 {
	// Implement logic to retrieve the IDs of the user's favorite assets
	// You can access the user's ID through u.ID
	return []int64{}
}


func (u *User) UpdateName(newName string) {
	u.Name = newName
}

func (u *User) GetAge() int {
	return 0
}

func (u *User) IsMale() bool {
	return false
}

// GetSocialMediaUsage return the duration of social media usage
func (u *User) GetSocialMediaUsage() time.Duration {	
	return 0
}

// GetLastMonthPurchaseCount returns the number of purchases
func (u *User) GetLastMonthPurchaseCount() int {
	return 0
}

// GDPR encrypt sensitive data
func (u *User) GDPR(ch CryptoHasher) (*User, error) {
	if u.Email == "" {
		return nil, errors.New("can't decrypt empty `Email` field")
	}
	encdata, err := ch.DecryptString(u.Email)
	if err != nil {
		return nil, err
	}
	u.Email = encdata

	return u, nil
}
