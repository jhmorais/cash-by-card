package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/cardMachine"
)

type DeleteCardMachineUseCase interface {
	Execute(ctx context.Context, cardMachineID int) (*output.DeleteCardMachine, error)
}
