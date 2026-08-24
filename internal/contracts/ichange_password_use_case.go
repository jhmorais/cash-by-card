package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
)

type ChangePasswordUseCase interface {
	Execute(ctx context.Context, email string, changePassword *input.ChangePassword) error
}
