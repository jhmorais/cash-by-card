package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
)

type FindPartnerByIDUseCase interface {
	Execute(ctx context.Context, clientID int) (*output.FindPartner, error)
}
