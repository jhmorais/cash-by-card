package input

type CreateCardMachine struct {
	Brand         []string `json:"brand"`
	PresentialTax float64  `json:"presentialTax"`
	OnlineTax     float64  `json:"onlineTax"`
	Installments  int      `json:"installments"`
}
