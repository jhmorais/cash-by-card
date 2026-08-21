package loan

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/loan"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/loan"
)

const (
	defaultPage  = 1
	defaultLimit = 10
	maxLimit     = 100
)

type listLoanUseCase struct {
	loanRepository repositories.LoanRepository
}

func NewListLoansUseCase(loanRepository repositories.LoanRepository) contracts.ListLoanUseCase {

	return &listLoanUseCase{
		loanRepository: loanRepository,
	}
}

func (l *listLoanUseCase) Execute(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) (*output.ListLoan, error) {
	if pagination == nil {
		pagination = &input.Pagination{Page: defaultPage, Limit: defaultLimit}
	}
	if pagination.Page < 1 {
		pagination.Page = defaultPage
	}
	if pagination.Limit < 1 {
		pagination.Limit = defaultLimit
	}
	if pagination.Limit > maxLimit {
		pagination.Limit = maxLimit
	}
	if filter == nil {
		filter = &input.ListLoanFilter{}
	}

	loans, total, err := l.loanRepository.ListLoan(ctx, filter, pagination)
	if err != nil {
		return nil, fmt.Errorf("error when list Loans on database: %v", err)
	}

	totalPages := 0
	if pagination.Limit > 0 {
		totalPages = int((total + int64(pagination.Limit) - 1) / int64(pagination.Limit))
	}

	if loans == nil {
		loans = []*entities.Loan{}
	}

	return &output.ListLoan{
		Loans:      loans,
		Total:      total,
		Page:       pagination.Page,
		Limit:      pagination.Limit,
		TotalPages: totalPages,
	}, nil
}
