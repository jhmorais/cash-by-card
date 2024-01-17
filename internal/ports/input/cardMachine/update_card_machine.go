package input

type UpdateCardMachine struct {
	ID            int      `json:"id,omitempty"`
	Brand         []string `json:"brand"`
	PresentialTax float64  `json:"presentialTax"`
	OnlineTax     float64  `json:"onlineTax"`
	Installments  int      `json:"installments"`
}
