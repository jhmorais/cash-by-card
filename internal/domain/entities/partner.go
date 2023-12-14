package entities

import "time"

type Partner struct {
	ID        int64  `gorm:"id" json:"id"`
	Name      string `gorm:"index" json:"name"`
	CPF       string `gorm:"index" json:"cpf"`
	PixKey    string `json:"pixKey"`
	Phone     string `json:"phone"`
	Endereco  string `json:"endereco"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
