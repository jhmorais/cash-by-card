package loan

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/loan"
)

type deleteLoanUseCase struct {
	loanRepository repositories.LoanRepository
}

func NewDeleteLoanUseCase(loanRepository repositories.LoanRepository) contracts.DeleteLoanUseCase {

	return &deleteLoanUseCase{
		loanRepository: loanRepository,
	}
}

func (c *deleteLoanUseCase) Execute(ctx context.Context, loanID int) (*output.DeleteLoan, error) {

	loanEntity, err := c.loanRepository.FindLoanByID(ctx, loanID)
	if err != nil {
		return nil, fmt.Errorf("failed to find loan '%d' at database: '%v'", loanID, err)
	}

	if loanEntity == nil || loanEntity.ID == 0 {
		return nil, fmt.Errorf("loanID not found")
	}

	err = c.loanRepository.DeleteLoan(ctx, loanEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to delete loan '%d'", loanEntity.ID)
	}

	output := &output.DeleteLoan{
		LoanID: loanEntity.ID,
	}

	return output, nil
}
