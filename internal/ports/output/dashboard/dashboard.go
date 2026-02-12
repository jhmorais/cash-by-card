package output

type Dashboard struct {
	TotalLoans    int          `json:"totalLoans"`
	TotalValue    float64      `json:"totalValue"`
	GrossProfit   float64      `json:"grossProfit"`
	Profit        float64      `json:"profit"`
	MachineAmount float64      `json:"machineAmount"`
	PartnerProfit float64      `json:"partnerProfit"`
	MonthlyLoans  MonthlyLoans `json:"monthlyLoans"`
}

type PartnerOutput struct {
	Name   string  `json:"name"`
	Qtt    int     `json:"qtt"`
	Profit float64 `json:"profit"`
}

type MonthlyLoans struct {
	Labels        []string  `json:"labels"`
	Total         []float64 `json:"total"`
	Gross         []float64 `json:"gross"`
	Net           []float64 `json:"net"`
	PartnerProfit []float64 `json:"partnerProfit"`
	MachineAmount []float64 `json:"machineAmount"`
}

type BestPartner struct {
	Partner string  `json:"partner"`
	Qtt     int     `json:"qtt"`
	Profit  float64 `json:"profit"`
}

type DashboardResponse struct {
	Dashboard    Dashboard     `json:"dashboard"`
	BestPartners []BestPartner `json:"bestPartners"`
}
