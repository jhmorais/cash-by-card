package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/user"
)

type DeleteUserUseCase interface {
	Execute(ctx context.Context, requesterEmail string, id int) (*output.DeleteUser, error)
}
