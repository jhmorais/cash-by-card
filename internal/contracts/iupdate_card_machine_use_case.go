package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/cardMachine"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/cardMachine"
)

type UpdateCardMachineUseCase interface {
	Execute(ctx context.Context, updateCardMachine *input.UpdateCardMachine) (*output.CreateCardMachine, error)
}
