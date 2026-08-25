package output

import "time"

type PartnerReportSummary struct {
	TotalLoans      int     `json:"totalLoans"`
	TotalCommission float64 `json:"totalCommission"`
}

type PartnerMonth struct {
	Month      int     `json:"month"` // 1-12
	Loans      int     `json:"loans"`
	Commission float64 `json:"commission"`
}

type PartnerYear struct {
	Year   int            `json:"year"`
	Months []PartnerMonth `json:"months"` // sempre 12 posições
}

type PartnerMonthDetail struct {
	LoanID     int       `json:"loanId"`
	Commission float64   `json:"commission"`
	CreatedAt  time.Time `json:"createdAt"`
	ClientName string    `json:"clientName"`
}

type PartnerReport struct {
	Summary      PartnerReportSummary `json:"summary"`
	Annual       []PartnerYear        `json:"annual"`
	CurrentMonth []PartnerMonthDetail `json:"currentMonth"`
	GeneratedAt  time.Time            `json:"generatedAt"`
}
