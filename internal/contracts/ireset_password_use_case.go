package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
)

type ResetPasswordUseCase interface {
	Execute(ctx context.Context, resetPassword *input.ResetPassword) error
}
