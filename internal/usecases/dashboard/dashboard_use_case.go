package dashboard

import (
	"context"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/dashboard"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/loan"
)

type dashboardUseCase struct {
	loanRepository repositories.LoanRepository
}

func NewDashboardUseCase(loanRepository repositories.LoanRepository) contracts.DashboardUseCase {

	return &dashboardUseCase{
		loanRepository: loanRepository,
	}
}

func (u *dashboardUseCase) Execute(ctx context.Context, month int, year int) (*output.DashboardResponse, error) {
	currentYear := time.Now().Year()

	totals, err := u.loanRepository.GetTotals(ctx, month, currentYear)
	if err != nil {
		return nil, fmt.Errorf("error fetching totals: %v", err)
	}

	partners, err := u.loanRepository.GetBestPartners(ctx, month, currentYear)
	if err != nil {
		return nil, fmt.Errorf("error fetching best partners: %v", err)
	}

	monthlyLoans, err := u.loanRepository.GetMonthlyLoans(ctx, year)
	if err != nil {
		return nil, fmt.Errorf("error fetching monthly loans: %v", err)
	}

	response := &output.DashboardResponse{
		Dashboard: output.Dashboard{
			TotalLoans:    totals.TotalLoans,
			TotalValue:    totals.TotalValue,
			GrossProfit:   totals.GrossProfit,
			Profit:        totals.Profit,
			PartnerProfit: totals.PartnerProfit,
			MachineAmount: totals.MachineAmount,
			MonthlyLoans:  *monthlyLoans,
		},
		BestPartners: partners,
	}

	return response, nil
}
