package contracts

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/ports/output"
)

type ListClientUseCase interface {
	Execute(ctx context.Context) (*output.ListClient, error)
}
