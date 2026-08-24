package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/loan"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
)

type ListLoanUseCase interface {
	Execute(ctx context.Context, filter *input.ListLoanFilter, pagination *input.Pagination) (*output.ListLoan, error)
}
