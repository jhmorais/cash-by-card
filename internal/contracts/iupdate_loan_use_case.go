package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/loan"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
)

type UpdateLoanUseCase interface {
	Execute(ctx context.Context, updateLoan *input.UpdateLoan) (*output.CreateLoan, error)
}
