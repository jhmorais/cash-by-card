package contracts

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/ports/output"
)

type FindClientUseCase interface {
	Execute(ctx context.Context, partnerID int, name string) (*output.FindClient, error)
}
