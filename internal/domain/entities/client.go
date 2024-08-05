package entities

import "time"

type Client struct {
	ID        int    `gorm:"id" json:"id"`
	Name      string `gorm:"size:250" json:"name"`
	PixType   int    `json:"pixType"`
	PixKey    string `json:"pixKey"`
	Phone     string `gorm:"size:15" json:"phone"`
	CPF       string `gorm:"size:15" json:"cpf"`
	PartnerID *int   `gorm:"index" json:"partnerId"`
	Documents string `json:"documents"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Partner   Partner `gorm:"foreignKey:PartnerID" json:"partner"`
}
