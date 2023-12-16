package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
)

type ListPartnerUseCase interface {
	Execute(ctx context.Context) (*output.ListPartner, error)
}
