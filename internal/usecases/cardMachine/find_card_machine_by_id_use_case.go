package cardMachine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/cardMachine"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/cardMachine"
)

type findCardMachineByIDUseCase struct {
	cardMachineRepository repositories.CardMachineRepository
}

func NewFindCardMachineByIDUseCase(cardMachineRepository repositories.CardMachineRepository) contracts.FindCardMachineByIDUseCase {

	return &findCardMachineByIDUseCase{
		cardMachineRepository: cardMachineRepository,
	}
}

func (c *findCardMachineByIDUseCase) Execute(ctx context.Context, cardMachineID int) (*output.FindCardMachine, error) {

	cardMachineEntity, err := c.cardMachineRepository.FindCardMachineByID(ctx, cardMachineID)
	if err != nil {
		return nil, fmt.Errorf("erro to find cardMachine '%d' at database: '%v'", cardMachineID, err)
	}

	if cardMachineEntity == nil || cardMachineEntity.ID == 0 {
		return nil, fmt.Errorf("cardMachineID not found")
	}

	var onlineTax map[string]interface{}
	err = json.Unmarshal(cardMachineEntity.OnlineTax, &onlineTax)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal online tax")
	}

	var presentialTax map[string]interface{}
	err = json.Unmarshal(cardMachineEntity.PresentialTax, &presentialTax)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal presential tax")
	}

	cardMachineOutput := &output.FindCardMachineOutput{
		ID:            cardMachineEntity.ID,
		Brand:         cardMachineEntity.Brand,
		Name:          cardMachineEntity.Name,
		PresentialTax: presentialTax,
		OnlineTax:     onlineTax,
		Installments:  cardMachineEntity.Installments,
		CreatedAt:     cardMachineEntity.CreatedAt,
		UpdatedAt:     cardMachineEntity.UpdatedAt,
	}

	output := &output.FindCardMachine{
		CardMachine: cardMachineOutput, //fazer o parse do byte[] para string
	}

	return output, nil
}
