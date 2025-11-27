package contracts

import (
	"context"

	output "github.com/jhmorais/cash-by-card/internal/ports/output/dashboard"
)

type DashboardUseCase interface {
	Execute(ctx context.Context, month int, year int) (*output.DashboardResponse, error)
}
