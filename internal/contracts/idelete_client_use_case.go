package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/client"
)

type DeleteClientUseCase interface {
	Execute(ctx context.Context, clientID int) (*output.DeleteClient, error)
}
