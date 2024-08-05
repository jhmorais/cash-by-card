package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
)

type FindPartnerByEmailUseCase interface {
	Execute(ctx context.Context, email string) (*output.ListPartner, error)
}
