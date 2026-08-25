package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
)

type PartnerReportUseCase interface {
	Execute(ctx context.Context, partnerUserEmail string) (*output.PartnerReport, error)
}
