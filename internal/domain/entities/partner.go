package entities

import "time"

type Partner struct {
	ID        int64  `gorm:"id"`
	Name      string `gorm:"index"`
	CPF       string `gorm:"index"`
	PixKey    string
	Phone     string
	Endereco  string
	CreatedAt time.Time
	UpdatedAt time.Time
}
