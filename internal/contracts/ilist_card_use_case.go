package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/card"
)

type ListCardUseCase interface {
	Execute(ctx context.Context) (*output.ListCard, error)
}
