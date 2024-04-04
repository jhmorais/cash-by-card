package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/card"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/card"
)

type CreateCardUseCase interface {
	Execute(ctx context.Context, createCard *[]input.CreateCard) ([]*output.CreateCard, error)
}
