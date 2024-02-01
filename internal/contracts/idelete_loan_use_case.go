package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
)

type DeleteLoanUseCase interface {
	Execute(ctx context.Context, lonaID int) (*output.DeleteLoan, error)
}
