package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
)

type FindLoanByClientIDUseCase interface {
	Execute(ctx context.Context, ClientID int) (*output.ListLoan, error)
}
