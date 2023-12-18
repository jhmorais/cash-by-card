package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/partner"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
)

type CreatePartnerUseCase interface {
	Execute(ctx context.Context, createPartner *input.CreatePartner) (*output.CreatePartner, error)
}
