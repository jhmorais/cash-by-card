package entities

import "time"

type CardMachine struct {
	ID            int    `gorm:"id" json:"id"`
	Brand         string `gorm:"size:250" json:"brand"`
	Name          string `gorm:"size:250" json:"name"`
	PresentialTax []byte `json:"presentialTax"` // quando for devolver para o front fazer um parse para string, JHSON.parse no front para acessar ao json
	OnlineTax     []byte `json:"onlineTax"`     // quando for devolver para o front fazer um parse para string, JHSON.parse no front para acessar ao json
	Installments  int    `gorm:"index" json:"installments"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
