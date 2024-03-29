package entities

import "time"

type Card struct {
	ID                int     `gorm:"id" json:"id"`
	PaymentType       string  `gorm:"size:250" json:"paymentType"`
	Value             float64 `json:"value"`
	Brand             string  `json:"brand"`
	Installments      int     `json:"installments"`
	InstallmentsValue float64 `json:"installmentsValue"`
	CardMachineName   string  `json:"cardMachineName"`
	LoanID            int     `gorm:"index" json:"loanId"`
	CardMachineID     int     `gorm:"index" json:"cardMachineId"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Loan              Loan        `gorm:"foreignKey:LoanID"`
	CardMachine       CardMachine `gorm:"foreignKey:CardMachineID"`
}
