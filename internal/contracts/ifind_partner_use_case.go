package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
)

type FindPartnerUseCase interface {
	Execute(ctx context.Context, partnerID int, name string) (*output.FindPartner, error)
}
