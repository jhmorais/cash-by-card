package card

import (
	"context"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/card"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/card"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/card"
)

type updateCardUseCase struct {
	cardRepository repositories.CardRepository
}

func NewUpdateCardUseCase(cardRepository repositories.CardRepository) contracts.UpdateCardUseCase {

	return &updateCardUseCase{
		cardRepository: cardRepository,
	}
}

func (c *updateCardUseCase) Execute(ctx context.Context, updateCard *input.UpdateCard) (*output.CreateCard, error) {

	if !(updateCard.PaymentType == "online" || updateCard.PaymentType == "present") {
		return nil, fmt.Errorf("failed updated card- PaymentType is empty")
	}

	if updateCard.Value == 0 {
		return nil, fmt.Errorf("failed updated card- value is empty")
	}

	if updateCard.Brand == "" {
		return nil, fmt.Errorf("failed updated card- brand is empty")
	}

	if updateCard.Installments <= 0 {
		return nil, fmt.Errorf("failed updated card- invalid number of installments")
	}
	if updateCard.LoanID <= 0 {
		return nil, fmt.Errorf("failed updated card- invalid LoanID")
	}
	if updateCard.CardMachineID <= 0 {
		return nil, fmt.Errorf("failed updated card- invalid CardMachineID")
	}

	cardEntity := &entities.Card{
		ID:            updateCard.ID,
		PaymentType:   updateCard.PaymentType,
		Value:         updateCard.Value,
		Brand:         updateCard.Brand,
		LoanID:        updateCard.LoanID,
		CardMachineID: updateCard.CardMachineID,
		CreatedAt:     time.Now(),
	}

	errUpdate := c.cardRepository.UpdateCard(ctx, cardEntity)
	if errUpdate != nil {
		return nil, fmt.Errorf("cannot update card at database: %v", errUpdate)
	}

	createCardOutput := &output.CreateCard{
		CardID: cardEntity.ID,
	}

	return createCardOutput, nil
}
