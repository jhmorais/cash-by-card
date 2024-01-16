package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/client"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/client"
)

type CreateClientUseCase interface {
	Execute(ctx context.Context, createClient *input.CreateClient) (*output.CreateClient, error)
}
