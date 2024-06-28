package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
)

type ListUserUseCase interface {
	Execute(ctx context.Context) ([]*output.FindUser, error)
}
