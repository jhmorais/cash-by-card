package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
)

type UpdateUserUseCase interface {
	Execute(ctx context.Context, requesterEmail string, updateUser *input.UpdateUser) (*output.UpdateUser, error)
}
