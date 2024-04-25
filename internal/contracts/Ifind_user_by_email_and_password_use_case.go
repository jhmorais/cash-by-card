package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
)

type FindUserByEmailAndPasswordUseCase interface {
	Execute(ctx context.Context, email string, password string) (*output.FindUser, error)
}
