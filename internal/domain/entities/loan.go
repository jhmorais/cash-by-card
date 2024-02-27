package entities

import "time"

type Loan struct {
	ID               int       `gorm:"id" json:"id"`
	ClientID         int       `gorm:"index" json:"clientId"`
	AskValue         float64   `json:"askValue"`
	Amount           float64   `json:"amount"`
	OperationPercent float64   `json:"operationPercent"`
	NumberCards      int       `json:"numberCards"`
	Cards            []Card    `json:"cards"`
	PartnerID        int       `gorm:"index" json:"partnerId"`
	GrossProfit      float64   `json:"grossProfit"`
	PartnerPercent   float64   `json:"partnerPercent"`
	PartnerAmount    float64   `json:"partnerAmount"`
	Profit           float64   `json:"profit"`
	PaymentStatus    string    `json:"paymentStatus"`
	ClientAmount     float64   `json:"clientAmount"`
	Type             int       `json:"type"`
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
	Partner          Partner   `gorm:"foreignKey:PartnerID" json:"partner"`
	Client           Client    `gorm:"foreignKey:ClientID" json:"client"`
}
