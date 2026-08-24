package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
)

type GetUserUseCase interface {
	Execute(ctx context.Context, email string) (*output.GetUser, error)
}
