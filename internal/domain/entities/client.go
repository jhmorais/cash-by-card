package entities

import "time"

type Client struct {
	ID        int    `gorm:"id" json:"id"`
	Name      string `gorm:"size:250" json:"name"`
	PixType   int    `json:"pixType"`
	PixKey    string `json:"pixKey"`
	PartnerID int    `gorm:"index" json:"partnerId"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Partner   Partner `gorm:"foreignKey:PartnerID" json:"partner"`
}
