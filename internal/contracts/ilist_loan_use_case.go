package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
)

type ListLoanUseCase interface {
	Execute(ctx context.Context) (*output.ListLoan, error)
}
