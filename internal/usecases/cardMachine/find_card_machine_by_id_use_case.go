package cardMachine

import (
	"context"
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

	output := &output.FindCardMachine{
		CardMachine: cardMachineEntity,
	}

	return output, nil
}
