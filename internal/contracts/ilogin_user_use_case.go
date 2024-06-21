package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
)

type LoginUseCase interface {
	Execute(ctx context.Context, userLogin *input.UserLogin) (string, error)
}
