package contracts

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/ports/input"
	"github.com/jhmorais/cash-by-card/internal/ports/output"
)

type CreateClientUseCase interface {
	Execute(ctx context.Context, createClient *input.CreateClient) (*output.CreateClient, error)
}
