package contracts

import (
	"context"

	input "github.com/jhmorais/cash-by-card/internal/ports/input/client"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
)

// PartnerClientsUseCase: CRUD de clientes escopado ao parceiro logado.
type PartnerClientsUseCase interface {
	ListClients(ctx context.Context, partnerUserEmail string) (*output.PartnerClients, error)
	CreateClient(ctx context.Context, partnerUserEmail string, createClient *input.CreateClient) (*output.PartnerCreateClient, error)
	UpdateClient(ctx context.Context, partnerUserEmail string, updateClient *input.UpdateClient) (*output.PartnerUpdateClient, error)
}
