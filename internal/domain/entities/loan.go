package entities

import "time"

type Loan struct {
	ID            int64 `gorm:"id"`
	ClientID      int64 `gorm:"index"`
	AskValue      float64
	Amount        float64
	NumberCards   int
	CardMachineID int64
	PartnerID     int `gorm:"index"`
	GrossValue    float64
	PartnerAmount float64
	Profit        float64
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Partner       Partner     `gorm:"foreignKey:PartnerID"`
	CardMachine   CardMachine `gorm:"foreignKey:CardMachineID"`
	Client        Client      `gorm:"foreignKey:ClientID"`
}
