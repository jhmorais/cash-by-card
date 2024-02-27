package input

type CreateCard struct {
	PaymentType       string  `json:"paymentType"`
	Value             float64 `json:"value"`
	Brand             string  `json:"brand"`
	Installments      int     `json:"installments"`
	InstallmentsValue float64 `json:"installmentsValue"`
	CardMachineName   string  `json:"cardMachineName"`
	LoanID            int     `json:"loanId"`
	CardMachineID     int     `json:"cardMachineId"`
}
