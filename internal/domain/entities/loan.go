package entities

import "time"

type Loan struct {
	ID               int     `gorm:"id" json:"id"`
	ClientID         int     `gorm:"index" json:"ClientId"`
	AskValue         float64 `json:"askValue"`
	Amount           float64 `json:"amount"`
	OperationPercent float64 `json:"operationPercent"`
	NumberCards      int     `json:"NumberCards"`
	Cards            []Card  `json:"cards"`
	PartnerID        int     `gorm:"index" json:"partnerId"`
	GrossProfit      float64 `json:"grossProfit"`
	PartnerPercent   float64 `json:"partnerPercent"`
	PartnerAmount    float64 `json:"partnerAmount"`
	Profit           float64 `json:"profit"`
	PaymentStatus    string  `json:"paymentStatus"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
	Partner          Partner `gorm:"foreignKey:PartnerID"`
	Client           Client  `gorm:"foreignKey:ClientID"`
}
