package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/cardMachine"
)

type FindCardMachineByIDUseCase interface {
	Execute(ctx context.Context, cardMachineID int) (*output.FindCardMachine, error)
}
