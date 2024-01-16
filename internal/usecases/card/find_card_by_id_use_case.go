package card

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/card"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/card"
)

type findCardByIDUseCase struct {
	cardRepository repositories.CardRepository
}

func NewFindCardByIDUseCase(cardRepository repositories.CardRepository) contracts.FindCardByIDUseCase {

	return &findCardByIDUseCase{
		cardRepository: cardRepository,
	}
}

func (c *findCardByIDUseCase) Execute(ctx context.Context, cardID int) (*output.FindCard, error) {

	cardEntity, err := c.cardRepository.FindCardByID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("erro to find client '%d' at database: '%v'", cardID, err)
	}

	if cardEntity == nil || cardEntity.ID == 0 {
		return nil, fmt.Errorf("cardID not found")
	}

	output := &output.FindCard{
		Card: cardEntity,
	}

	return output, nil
}
