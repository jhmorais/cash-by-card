package card

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/card"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/card"
)

type listCardUseCase struct {
	cardRepository repositories.CardRepository
}

func NewListCardsUseCase(cardRepository repositories.CardRepository) contracts.ListCardUseCase {

	return &listCardUseCase{
		cardRepository: cardRepository,
	}
}

func (l *listCardUseCase) Execute(ctx context.Context) (*output.ListCard, error) {
	var err error
	output := &output.ListCard{Cards: []*entities.Card{}}

	output.Cards, err = l.cardRepository.ListCard(ctx)
	if err != nil {
		return nil, fmt.Errorf("error when list cards on database: %v", err)
	}

	return output, nil
}
