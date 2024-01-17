package cardMachine

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
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
	output := &output.ListCardMachine{CardMachines: []*entities.CardMachine{}}

	output.CardMachines, err = l.cardMachineRepository.ListCardMachine(ctx)
	if err != nil {
		return nil, fmt.Errorf("error when list cardMachines on database: %v", err)
	}

	return output, nil
}
