package contracts

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/ports/output"
)

type FindClientByIDUseCase interface {
	Execute(ctx context.Context, clientID int) (*output.FindClient, error)
}
