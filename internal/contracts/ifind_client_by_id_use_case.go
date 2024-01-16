package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/client"
)

type FindClientByIDUseCase interface {
	Execute(ctx context.Context, clientID int) (*output.FindClient, error)
}
