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

type createCardMachineUseCase struct {
	cardMachineRepository repositories.CardMachineRepository
}

func NewCreateCardMachineUseCase(cardMachineRepository repositories.CardMachineRepository) contracts.CreateCardMachineUseCase {

	return &createCardMachineUseCase{
		cardMachineRepository: cardMachineRepository,
	}
}

func (c *createCardMachineUseCase) Execute(ctx context.Context, createCardMachine *input.CreateCardMachine) (*output.CreateCardMachine, error) {

	if createCardMachine.Name == "" {
		return nil, fmt.Errorf("não é possível criar a maquininha sem nome")
	}

	if createCardMachine.Installments <= 0 {
		return nil, fmt.Errorf("não é possível criar a maquininha sem número de parcelas")
	}
	if len(createCardMachine.Brand) == 0 {
		return nil, fmt.Errorf("não é possível criar a maquininha sem bandeira")
	}

	if createCardMachine.OnlineTax == nil {
		return nil, fmt.Errorf("cannot create cardMachine without valid OnlineTax")
	}

	if createCardMachine.PresentialTax == nil {
		return nil, fmt.Errorf("cannot create cardMachine without valid PresentialTax")
	}

	bOnlineTax, err := json.Marshal(createCardMachine.OnlineTax)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal online tax")
	}

	bPresentialTax, err := json.Marshal(createCardMachine.PresentialTax)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal presential tax")
	}

	brands, err := json.Marshal(createCardMachine.Brand)
	if err != nil {
		return nil, err
	}
	// createdAt := time.Now()
	cardMachineEntity := &entities.CardMachine{
		Brand:         string(brands),
		Name:          createCardMachine.Name,
		Installments:  createCardMachine.Installments,
		OnlineTax:     bOnlineTax,
		PresentialTax: bPresentialTax,
		CreatedAt:     time.Now(),
	}

	err = c.cardMachineRepository.CreateCardMachine(ctx, cardMachineEntity)
	if err != nil {
		return nil, fmt.Errorf("não foi possível salvar a maquininha: %v", err)
	}

	createCardMachineOutput := &output.CreateCardMachine{
		CardMachineID: cardMachineEntity.ID,
	}

	return createCardMachineOutput, nil
}
