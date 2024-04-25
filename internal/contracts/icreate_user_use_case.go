package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/user"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
)

type CreateUserUseCase interface {
	Execute(ctx context.Context, createUser *input.CreateUser) (*output.CreateUser, error)
}
