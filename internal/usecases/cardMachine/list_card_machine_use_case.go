package cardMachine

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/cardMachine"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/cardMachine"
)

type listCardMachineUseCase struct {
	cardMachineRepository repositories.CardMachineRepository
}

func NewListCardMachinesUseCase(cardMachineRepository repositories.CardMachineRepository) contracts.ListCardMachineUseCase {

	return &listCardMachineUseCase{
		cardMachineRepository: cardMachineRepository,
	}
}

func (l *listCardMachineUseCase) Execute(ctx context.Context) (*output.ListCardMachine, error) {
	var err error
	result := &output.ListCardMachine{CardMachines: []*output.FindCardMachineOutput{}}

	cardMachinesEntity, err := l.cardMachineRepository.ListCardMachine(ctx)
	if err != nil {
		return nil, fmt.Errorf("error when list cardMachines on database: %v", err)
	}

	for _, cardMachine := range cardMachinesEntity {
		var onlineTax map[string]interface{}
		err = json.Unmarshal(cardMachine.OnlineTax, &onlineTax)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal online tax")
		}

		var presentialTax map[string]interface{}
		err = json.Unmarshal(cardMachine.PresentialTax, &presentialTax)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal presential tax")
		}

		result.CardMachines = append(result.CardMachines, &output.FindCardMachineOutput{
			ID:            cardMachine.ID,
			Brand:         cardMachine.Brand,
			Name:          cardMachine.Name,
			PresentialTax: presentialTax,
			OnlineTax:     onlineTax,
			Installments:  cardMachine.Installments,
			CreatedAt:     cardMachine.CreatedAt,
			UpdatedAt:     cardMachine.UpdatedAt,
		})
	}

	return result, nil
}
