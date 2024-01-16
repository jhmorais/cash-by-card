package card

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/card"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/card"
)

type deleteCardUseCase struct {
	cardRepository repositories.CardRepository
}

func NewDeleteCardUseCase(cardRepository repositories.CardRepository) contracts.DeleteCardUseCase {

	return &deleteCardUseCase{
		cardRepository: cardRepository,
	}
}

func (c *deleteCardUseCase) Execute(ctx context.Context, cardID int) (*output.DeleteCard, error) {

	cardEntity, err := c.cardRepository.FindCardByID(ctx, cardID)
	if err != nil {
		return nil, fmt.Errorf("failed to find card '%d' at database: '%v'", cardID, err)
	}

	if cardEntity == nil || cardEntity.ID == 0 {
		return nil, fmt.Errorf("cardID not found")
	}

	err = c.cardRepository.DeleteCard(ctx, cardEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to delete card '%d'", cardEntity.ID)
	}

	output := &output.DeleteCard{
		CardID: cardEntity.ID,
	}

	return output, nil
}
