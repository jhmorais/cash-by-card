package entities

import "time"

type Card struct {
	ID                int     `gorm:"id" json:"id"`
	PaymentType       string  `gorm:"size:250" json:"paymentType"`
	Value             float64 `json:"value"`
	Brand             string  `json:"brand"`
	Installments      int     `json:"installments"`
	InstallmentsValue float64 `json:"installmentsValue"`
	MachineValue      float64 `json:"machineValue"`
	CardMachineName   string  `json:"cardMachineName"`
	ClientAmount      float64 `json:"clientAmount"`
	GrossProfit       float64 `json:"grossProfit"`
	LoanID            int     `gorm:"index" json:"loanId"`
	CardMachineID     int     `gorm:"index" json:"cardMachineId"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	Loan              Loan        `gorm:"foreignKey:LoanID"`
	CardMachine       CardMachine `gorm:"foreignKey:CardMachineID"`
}
