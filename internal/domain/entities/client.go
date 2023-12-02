package entities

import "time"

type Client struct {
	ID        int    `gorm:"id"`
	Name      string `gorm:"size:250"`
	PixType   int
	PixKey    string
	PartnerID int `gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Partner   Partner `gorm:"foreignKey:PartnerID"`
}
