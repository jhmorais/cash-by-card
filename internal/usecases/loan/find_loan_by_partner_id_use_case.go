package loan

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/loan"
)

type findLoanByPartnerIDUseCase struct {
	loanRepository repositories.LoanRepository
}

func NewFindLoansByPartnerIDUseCase(loanRepository repositories.LoanRepository) contracts.FindLoanByParnterIDUseCase {

	return &findLoanByPartnerIDUseCase{
		loanRepository: loanRepository,
	}
}

func (c *findLoanByPartnerIDUseCase) Execute(ctx context.Context, PartnerID int) (*output.ListLoan, error) {

	loansEntity, err := c.loanRepository.FindLoanByPartnerID(ctx, PartnerID)
	if err != nil {
		return nil, fmt.Errorf("erro to find loan with PartnerId: '%v' at database: '%v'", PartnerID, err)
	}

	if len(loansEntity) == 0 {
		return nil, fmt.Errorf("loans not found")
	}

	output := &output.ListLoan{
		Loans: loansEntity,
	}

	return output, nil
}
