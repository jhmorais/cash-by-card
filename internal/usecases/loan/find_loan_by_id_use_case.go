package loan

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/loan"
)

type findLoanByIDUseCase struct {
	loanRepository repositories.LoanRepository
}

func NewFindLoanByIDUseCase(loanRepository repositories.LoanRepository) contracts.FindLoanByIDUseCase {

	return &findLoanByIDUseCase{
		loanRepository: loanRepository,
	}
}

func (c *findLoanByIDUseCase) Execute(ctx context.Context, cardID int) (*output.FindLoan, error) {

	loanEntity, err := c.loanRepository.FindLoanByID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("erro to find loan '%d' at database: '%v'", cardID, err)
	}

	if loanEntity == nil || loanEntity.ID == 0 {
		return nil, fmt.Errorf("LoanID not found")
	}

	output := &output.FindLoan{
		Loan: loanEntity,
	}

	return output, nil
}
