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

func (c *updateCardUseCase) Execute(ctx context.Context, updateCards []input.UpdateCard) (*output.UpdateCards, error) {
	cardsEntity := []*entities.Card{}
	for _, card := range updateCards {
		if err := validateFields(card); err != nil {
			return nil, err
		}
		cardEntity := &entities.Card{
			ID:                card.ID,
			PaymentType:       card.PaymentType,
			Value:             card.Value,
			MachineValue:      card.MachineValue,
			Installments:      card.Installments,
			InstallmentsValue: card.InstallmentsValue,
			Brand:             card.Brand,
			LoanID:            card.LoanID,
			CardMachineID:     card.CardMachineID,
			CardMachineName:   card.CardMachineName,
			ClientAmount:      card.ClientAmount,
			GrossProfit:       card.GrossProfit,
			UpdatedAt:         time.Now(),
		}
		cardsEntity = append(cardsEntity, cardEntity)
	}

	errUpdate := c.cardRepository.UpdateCard(ctx, cardsEntity)
	if errUpdate != nil {
		return nil, fmt.Errorf("cannot update card at database: %v", errUpdate)
	}

	createCardOutput := &output.UpdateCards{}
	for _, card := range cardsEntity {
		createCardOutput.CardIDs = append(createCardOutput.CardIDs, card.ID)
	}

	return createCardOutput, nil
}

func validateFields(updateCard input.UpdateCard) error {
	if !(updateCard.PaymentType == "onlineTax" || updateCard.PaymentType == "presentialTax") {
		return fmt.Errorf("failed updated card- PaymentType is empty")
	}

	if updateCard.Value == 0 {
		return fmt.Errorf("failed updated card- value is empty")
	}

	if updateCard.Brand == "" {
		return fmt.Errorf("failed updated card- brand is empty")
	}

	if updateCard.Installments <= 0 {
		return fmt.Errorf("failed updated card- invalid number of installments")
	}
	if updateCard.LoanID <= 0 {
		return fmt.Errorf("failed updated card- invalid LoanID")
	}
	if updateCard.CardMachineID <= 0 {
		return fmt.Errorf("failed updated card- invalid CardMachineID")
	}
	if updateCard.CardMachineName == "" {
		return fmt.Errorf("cannot updated card without valid CardMachine name")
	}

	return nil
}
