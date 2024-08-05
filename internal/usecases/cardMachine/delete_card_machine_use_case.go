package cardMachine

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/cardMachine"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/cardMachine"
)

type deleteCardMachineUseCase struct {
	cardMachineRepository repositories.CardMachineRepository
}

func NewDeleteCardMachineUseCase(cardMachineRepository repositories.CardMachineRepository) contracts.DeleteCardMachineUseCase {

	return &deleteCardMachineUseCase{
		cardMachineRepository: cardMachineRepository,
	}
}

func (c *deleteCardMachineUseCase) Execute(ctx context.Context, cardMachineID int) (*output.DeleteCardMachine, error) {

	cardMachineEntity, err := c.cardMachineRepository.FindCardMachineByID(ctx, cardMachineID)
	if err != nil {
		return nil, fmt.Errorf("failed to find cardMachine '%d' at database: '%v'", cardMachineID, err)
	}

	if cardMachineEntity == nil || cardMachineEntity.ID == 0 {
		return nil, fmt.Errorf("cardMachineID not found")
	}

	err = c.cardMachineRepository.DeleteCardMachine(ctx, cardMachineEntity)
	if err != nil {
		return nil, fmt.Errorf("falha em deletar a maquininha '%d'", cardMachineEntity.ID)
	}

	output := &output.DeleteCardMachine{
		CardMachineID: cardMachineEntity.ID,
	}

	return output, nil
}
