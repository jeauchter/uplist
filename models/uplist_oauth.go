package models

import (
	"time"

	"gorm.io/gorm"
)

// UplistOauth struct
type UplistOauth struct {
	gorm.Model
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	ExpiresAt      time.Time `json:"expires_at"`
	ShopID         int       `json:"shop_id" gorm:"default:0"`
	ShopName       string    `json:"shop_name"`
	APIKey         string
	APISecret      string
	APIRedirectURL string
	CallbackURL    string
}
