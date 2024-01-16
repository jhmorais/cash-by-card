package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/card"
)

type FindCardByLoanIDUseCase interface {
	Execute(ctx context.Context, LoanID int) (*output.ListCard, error)
}
