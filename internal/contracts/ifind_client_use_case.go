package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/client"
)

type FindClientUseCase interface {
	Execute(ctx context.Context, partnerID int, name string) (*output.FindClient, error)
}
