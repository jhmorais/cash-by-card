package cardMachine

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/cardMachine"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/cardMachine"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/cardMachine"
)

type updateCardMachineUseCase struct {
	cardMachineRepository repositories.CardMachineRepository
}

func NewUpdateCardMachineUseCase(cardMachineRepository repositories.CardMachineRepository) contracts.UpdateCardMachineUseCase {

	return &updateCardMachineUseCase{
		cardMachineRepository: cardMachineRepository,
	}
}

func (c *updateCardMachineUseCase) Execute(ctx context.Context, updateCardMachine *input.UpdateCardMachine) (*output.CreateCardMachine, error) {

	if len(updateCardMachine.Brand) == 0 {
		return nil, fmt.Errorf("failed updated cardMachine- Brand is empty")
	}
	if updateCardMachine.Installments <= 0 {
		return nil, fmt.Errorf("failed updated cardMachine- invalid number of installments")
	}
	if updateCardMachine.OnlineTax <= 0 {
		return nil, fmt.Errorf("failed updated cardMachine- invalid OnlineTax")
	}
	if updateCardMachine.PresentialTax <= 0 {
		return nil, fmt.Errorf("failed updated cardMachine- invalid PresentialTax")
	}

	brands, err := json.Marshal(updateCardMachine.Brand)
	if err != nil {
		return nil, err
	}

	cardMachineEntity := &entities.CardMachine{
		ID:            updateCardMachine.ID,
		Brand:         string(brands),
		Installments:  updateCardMachine.Installments,
		OnlineTax:     updateCardMachine.OnlineTax,
		PresentialTax: updateCardMachine.PresentialTax,
		CreatedAt:     time.Now(),
	}

	errUpdate := c.cardMachineRepository.UpdateCardMachine(ctx, cardMachineEntity)
	if errUpdate != nil {
		return nil, fmt.Errorf("cannot update cardMchine at database: %v", errUpdate)
	}

	createCardMachineOutput := &output.CreateCardMachine{
		CardMachineID: cardMachineEntity.ID,
	}

	return createCardMachineOutput, nil
}
