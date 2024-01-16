package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/card"
)

type FindCardByIDUseCase interface {
	Execute(ctx context.Context, cardID int) (*output.FindCard, error)
}
