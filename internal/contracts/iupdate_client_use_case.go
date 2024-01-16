package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/client"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/client"
)

type UpdateClientUseCase interface {
	Execute(ctx context.Context, updateClient *input.UpdateClient) (*output.CreateClient, error)
}
