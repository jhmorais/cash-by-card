package loan

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/loan"
)

type listLoanUseCase struct {
	loanRepository repositories.LoanRepository
}

func NewListLoansUseCase(loanRepository repositories.LoanRepository) contracts.ListLoanUseCase {

	return &listLoanUseCase{
		loanRepository: loanRepository,
	}
}

func (l *listLoanUseCase) Execute(ctx context.Context) (*output.ListLoan, error) {
	var err error
	output := &output.ListLoan{Loans: []*entities.Loan{}}

	output.Loans, err = l.loanRepository.ListLoan(ctx)
	if err != nil {
		return nil, fmt.Errorf("error when list Loans on database: %v", err)
	}

	return output, nil
}
