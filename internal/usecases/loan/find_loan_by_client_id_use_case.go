package loan

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/loan"
)

type findLoanByClientIDUseCase struct {
	loanRepository repositories.LoanRepository
}

func NewFindLoanByClientIDUseCase(loanRepository repositories.LoanRepository) contracts.FindLoanByClientIDUseCase {

	return &findLoanByClientIDUseCase{
		loanRepository: loanRepository,
	}
}

func (c *findLoanByClientIDUseCase) Execute(ctx context.Context, ClientID int) (*output.ListLoan, error) {

	loansEntity, err := c.loanRepository.FindLoanByClientID(ctx, ClientID)
	if err != nil {
		return nil, fmt.Errorf("erro to find loan with ClientId: '%v' at database: '%v'", ClientID, err)
	}

	if len(loansEntity) == 0 {
		return nil, fmt.Errorf("loans not found")
	}

	output := &output.ListLoan{
		Loans: loansEntity,
	}

	return output, nil
}
