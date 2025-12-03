package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
)

type UpdateUserUseCase interface {
	Execute(ctx context.Context, updateClient *input.UpdateUser) error
}
