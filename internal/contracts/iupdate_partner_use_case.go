package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/partner"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
)

type UpdatePartnerUseCase interface {
	Execute(ctx context.Context, updatePartner *input.UpdatePartner) (*output.CreatePartner, error)
}
