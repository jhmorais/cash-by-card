package entities

import "time"

type CardMachine struct {
	ID            string  `gorm:"id"`
	Brand         string  `gorm:"size:250"`
	PresentialTax float64 `gorm:"index"`
	OnlineTax     float64 `gorm:"index"`
	Installments  int     `gorm:"index"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}
