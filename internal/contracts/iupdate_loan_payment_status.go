package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/loan"
)

type UpdateLoanPaymentStatusUseCase interface {
	Execute(ctx context.Context, UpdateLoanPaymentStatus *input.UpdateLoanPaymentStatus) error
}
