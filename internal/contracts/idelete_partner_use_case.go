package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
)

type DeletePartnerUseCase interface {
	Execute(ctx context.Context, clientID int) (*output.DeletePartner, error)
}
