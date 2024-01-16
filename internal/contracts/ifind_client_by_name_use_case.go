package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/client"
)

type FindClientByNameUseCase interface {
	Execute(ctx context.Context, name string) (*output.ListClient, error)
}
