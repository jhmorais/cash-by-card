package output

type Dashboard struct {
	TotalLoans   int          `json:"totalLoans"`
	TotalValue   float64      `json:"totalValue"`
	GrossProfit  float64      `json:"grossProfit"`
	Profit       float64      `json:"profit"`
	MonthlyLoans MonthlyLoans `json:"monthlyLoans"`
}

type PartnerOutput struct {
	Name string `json:"name"`
	Qtt  int    `json:"qtt"`
}

type MonthlyLoans struct {
	Labels []string  `json:"labels"`
	Total  []float64 `json:"total"`
	Gross  []float64 `json:"gross"`
	Net    []float64 `json:"net"`
}

type BestPartner struct {
	Partner string `json:"partner"`
	Qtt     int    `json:"qtt"`
}

type DashboardResponse struct {
	Dashboard    Dashboard     `json:"dashboard"`
	BestPartners []BestPartner `json:"bestPartners"`
}
