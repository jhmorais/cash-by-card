package entities

import "time"

type CardMachine struct {
	ID            int     `gorm:"id" json:"id"`
	Brand         string  `gorm:"size:250" json:"brand"`
	PresentialTax float64 `gorm:"index" json:"presentialTax"`
	OnlineTax     float64 `gorm:"index" json:"onlineTax"`
	Installments  int     `gorm:"index" json:"installments"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
