package contracts

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/ports/output"
)

type FindClientByNameUseCase interface {
	Execute(ctx context.Context, name string) (*output.ListClient, error)
}
