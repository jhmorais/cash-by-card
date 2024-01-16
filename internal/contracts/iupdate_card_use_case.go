package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/card"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/card"
)

type UpdateCardUseCase interface {
	Execute(ctx context.Context, updateCard *input.UpdateCard) (*output.CreateCard, error)
}
