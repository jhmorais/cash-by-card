package input

type Card struct {
	ID                int     `json:"id"`
	PaymentType       string  `json:"paymentType"`
	Value             float64 `json:"value"`
	Brand             string  `json:"brand"`
	Installments      int     `json:"installments"`
	InstallmentsValue float64 `json:"installmentsValue"`
	CardMachineName   string  `json:"cardMachineName"`
	GrossProfit       float64 `json:"grossProfit"`
	ClientAmount      float64 `json:"clientAmount"`
	MachineValue      float64 `json:"machineValue"`
	LoanID            int     `json:"loanId"`
	CardMachineID     int     `json:"cardMachineId"`
}
