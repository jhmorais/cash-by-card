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

type createCardUseCase struct {
	cardRepository repositories.CardRepository
}

func NewCreateCardUseCase(cardRepository repositories.CardRepository) contracts.CreateCardUseCase {

	return &createCardUseCase{
		cardRepository: cardRepository,
	}
}

func (c *createCardUseCase) Execute(ctx context.Context, createCard *input.CreateCard) (*output.CreateCard, error) {

	if !(createCard.PaymentType == "online" || createCard.PaymentType == "present") {
		return nil, fmt.Errorf("cannot create a card with invalid payment type")
	}

	if createCard.Value == 0 {
		return nil, fmt.Errorf("cannot create a card without value")
	}

	if createCard.Brand == "" {
		return nil, fmt.Errorf("cannot create a card without brand")
	}

	if createCard.Installments <= 0 {
		return nil, fmt.Errorf("cannot create card without valid number of installments")
	}
	if createCard.LoanID <= 0 {
		return nil, fmt.Errorf("cannot create card without valid LoanID")
	}
	if createCard.CardMachineID <= 0 {
		return nil, fmt.Errorf("cannot create card without valid CardMachineID")
	}

	cardEntity := &entities.Card{
		PaymentType:   createCard.PaymentType,
		Value:         createCard.Value,
		Brand:         createCard.Brand,
		Installments:  createCard.Installments,
		LoanID:        createCard.LoanID,
		CardMachineID: createCard.CardMachineID,
		CreatedAt:     time.Now(),
	}

	err := c.cardRepository.CreateCard(ctx, cardEntity)
	if err != nil {
		return nil, fmt.Errorf("cannot save card at database: %v", err)
	}

	createCardOutput := &output.CreateCard{
		CardID: cardEntity.ID,
	}

	return createCardOutput, nil
}
