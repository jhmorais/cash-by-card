package partner

import (
	"context"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
	repoLoan "github.com/jhmorais/cash-by-card/internal/repositories/loan"
	repoPartner "github.com/jhmorais/cash-by-card/internal/repositories/partner"
)

type partnerReportUseCase struct {
	loanRepository    repoLoan.LoanRepository
	partnerRepository repoPartner.PartnerRepository
}

func NewPartnerReportUseCase(loanRepository repoLoan.LoanRepository, partnerRepository repoPartner.PartnerRepository) contracts.PartnerReportUseCase {
	return &partnerReportUseCase{loanRepository: loanRepository, partnerRepository: partnerRepository}
}

func (p *partnerReportUseCase) Execute(ctx context.Context, partnerUserEmail string) (*output.PartnerReport, error) {
	report := &output.PartnerReport{
		Summary:      output.PartnerReportSummary{},
		Annual:       []output.PartnerYear{},
		CurrentMonth: []output.PartnerMonthDetail{},
		GeneratedAt:  time.Now(),
	}
	partners, err := p.partnerRepository.FindPartnerByEmail(ctx, partnerUserEmail)
	if err != nil {
		return nil, err
	}
	if len(partners) == 0 || partners[0].Email == "" {
		return report, nil // sem entidade: relatório vazio
	}
	loans, err := p.loanRepository.FindLoanByPartnerID(ctx, partners[0].ID)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	type key struct{ year, month int }
	agg := map[key]*output.PartnerMonth{}
	for _, l := range loans {
		created := l.CreatedAt
		report.Summary.TotalLoans++
		report.Summary.TotalCommission += l.PartnerAmount

		k := key{created.Year(), int(created.Month())}
		m, ok := agg[k]
		if !ok {
			m = &output.PartnerMonth{Month: k.month}
			agg[k] = m
		}
		m.Loans++
		m.Commission += l.PartnerAmount

		if created.Year() == now.Year() && int(created.Month()) == int(now.Month()) {
			report.CurrentMonth = append(report.CurrentMonth, output.PartnerMonthDetail{
				LoanID:     l.ID,
				Commission: l.PartnerAmount,
				CreatedAt:  created,
				ClientName: l.Client.Name,
			})
		}
	}

	if len(agg) > 0 {
		minYear := now.Year()
		for k := range agg {
			if k.year < minYear {
				minYear = k.year
			}
		}
		for y := minYear; y <= now.Year(); y++ {
			py := output.PartnerYear{Year: y, Months: make([]output.PartnerMonth, 12)}
			for m := 1; m <= 12; m++ {
				py.Months[m-1] = output.PartnerMonth{Month: m}
				if agg[key{y, m}] != nil {
					py.Months[m-1] = *agg[key{y, m}]
				}
			}
			report.Annual = append(report.Annual, py)
		}
	}
	return report, nil
}
