package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
)

type FindPartnerByNameUseCase interface {
	Execute(ctx context.Context, name string) (*output.ListPartner, error)
}
