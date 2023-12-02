package entities

import "time"

type Card struct {
	ID            string `gorm:"id"`
	PaymentType   string `gorm:"size:250"`
	Value         float64
	Brand         string
	Installments  int
	LoanID        int `gorm:"index"`
	CardMachineID int `gorm:"index"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Loan          Loan        `gorm:"foreignKey:LoanID"`
	CardMachine   CardMachine `gorm:"foreignKey:CardMachineID"`
}
