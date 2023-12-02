package contracts

import (
	"context"

	"github.com/jhmorais/cash-by-card/internal/ports/output"
)

type DeleteClientUseCase interface {
	Execute(ctx context.Context, deviceID string) (*output.DeleteClient, error)
}
