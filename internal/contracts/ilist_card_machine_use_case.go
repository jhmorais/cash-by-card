package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/cardMachine"
)

type ListCardMachineUseCase interface {
	Execute(ctx context.Context) (*output.ListCardMachine, error)
}
