package input

import input "github.com/jhmorais/cash-by-card/internal/ports/input/card"

type CreateLoan struct {
	ClientID         int          `json:"clientId"`
	AskValue         float64      `json:"askValue"`
	OperationPercent float64      `json:"operationPercent"`
	Amount           float64      `json:"amount"`
	NumberCards      int          `json:"numberCards"`
	Cards            []input.Card `json:"cards"`
	GrossProfit      float64      `json:"grossProfit"`
	PartnerID        int          `json:"partnerId"`
	PartnerPercent   float64      `json:"partnerPercent"`
	PartnerAmount    float64      `json:"partnerAmount"`
	Profit           float64      `json:"profit"`
	PaymentStatus    string       `json:"paymentStatus"`
}
