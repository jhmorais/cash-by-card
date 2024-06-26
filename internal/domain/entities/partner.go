package entities

import "time"

type Partner struct {
	ID        int    `gorm:"id" json:"id"`
	Name      string `gorm:"index" json:"name"`
	CPF       string `gorm:"index" json:"cpf"`
	PixKey    string `json:"pixKey"`
	Phone     string `json:"phone"`
	Address   string `json:"address"`
	PixType   int    `json:"pixType"`
	Email     string `json:"email"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
