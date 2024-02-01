package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/loan"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/loan"
)

type CreateLoanUseCase interface {
	Execute(ctx context.Context, createClient *input.CreateLoan) (*output.CreateLoan, error)
}
