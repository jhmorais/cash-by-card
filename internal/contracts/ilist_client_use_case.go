package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/client"
)

type ListClientUseCase interface {
	Execute(ctx context.Context) (*output.ListClient, error)
}
