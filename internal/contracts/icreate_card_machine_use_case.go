package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/cardMachine"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/cardMachine"
)

type CreateCardMachineUseCase interface {
	Execute(ctx context.Context, createCardMachine *input.CreateCardMachine) (*output.CreateCardMachine, error)
}
