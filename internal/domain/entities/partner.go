package entities

import "time"

type Partner struct {
	ID        int    `gorm:"id" json:"id"`
	Name      string `gorm:"index" json:"name"`
	CPF       string `gorm:"index" json:"cpf"`
	PixKey    string `json:"pixKey"`
	Telefone  string `json:"phone"`
	Endereco  string `json:"address"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
