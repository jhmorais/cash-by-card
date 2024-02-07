package input

type UpdateCardMachine struct {
	ID            int                    `json:"id,omitempty"`
	Brand         []string               `json:"brand"`
	Name          string                 `json:"name"`
	PresentialTax map[string]interface{} `json:"presentialTax"`
	OnlineTax     map[string]interface{} `json:"onlineTax"`
	Installments  int                    `json:"installments"`
}
