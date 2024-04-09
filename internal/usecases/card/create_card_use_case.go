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

func (c *createCardUseCase) Execute(ctx context.Context, createCard *[]input.CreateCard) ([]*output.CreateCard, error) {
	var createCardOutput []*output.CreateCard
	for _, createCardInput := range *createCard {

		if !(createCardInput.PaymentType == "onlineTax" || createCardInput.PaymentType == "presentialTax") {
			return nil, fmt.Errorf("cannot create a card with invalid payment type")
		}

		if createCardInput.Value == 0 {
			return nil, fmt.Errorf("cannot create a card without value")
		}

		if createCardInput.Brand == "" {
			return nil, fmt.Errorf("cannot create a card without brand")
		}

		if createCardInput.Installments <= 0 {
			return nil, fmt.Errorf("cannot create card without valid number of installments")
		}

		if createCardInput.MachineValue <= 0 {
			return nil, fmt.Errorf("cannot create card without valid number of machine value")
		}

		if int(createCardInput.InstallmentsValue) <= 0 {
			return nil, fmt.Errorf("cannot create card without valid number of installmentsValue")
		}

		if createCardInput.LoanID <= 0 {
			return nil, fmt.Errorf("cannot create card without valid LoanID")
		}
		if createCardInput.CardMachineID <= 0 {
			return nil, fmt.Errorf("cannot create card without valid CardMachineID")
		}
		if createCardInput.CardMachineName == "" {
			return nil, fmt.Errorf("cannot create card without valid CardMachine name")
		}

		cardEntity := &entities.Card{
			PaymentType:       createCardInput.PaymentType,
			Value:             createCardInput.Value,
			Brand:             createCardInput.Brand,
			Installments:      createCardInput.Installments,
			InstallmentsValue: createCardInput.InstallmentsValue,
			MachineValue:      createCardInput.MachineValue,
			LoanID:            createCardInput.LoanID,
			CardMachineID:     createCardInput.CardMachineID,
			CardMachineName:   createCardInput.CardMachineName,
			ClientAmount:      createCardInput.ClientAmount,
			GrossProfit:       createCardInput.GrossProfit,
			CreatedAt:         time.Now(),
		}

		err := c.cardRepository.CreateCard(ctx, cardEntity)
		if err != nil {
			return nil, fmt.Errorf("cannot save card at database: %v", err)
		}

		createCardOutput = append(createCardOutput, &output.CreateCard{
			CardID: cardEntity.ID,
		})

	}
	if len(createCardOutput) == 0 {
		return nil, fmt.Errorf("empty input provided")
	}
	return createCardOutput, nil
}
