package contracts

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/ports/input"
	"github.com/jhmorais/cash-by-card/internal/ports/output"
)

type UpdateClientUseCase interface {
	Execute(ctx context.Context, updateClient *input.UpdateClient) (*output.CreateClient, error)
}
