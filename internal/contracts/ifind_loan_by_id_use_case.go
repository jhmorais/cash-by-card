package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
)

type FindLoanByIDUseCase interface {
	Execute(ctx context.Context, LoanID int) (*output.FindLoan, error)
}
