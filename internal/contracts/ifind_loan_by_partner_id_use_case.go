package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
)

type FindLoanByParnterIDUseCase interface {
	Execute(ctx context.Context, ParnterID int) (*output.ListLoan, error)
}
