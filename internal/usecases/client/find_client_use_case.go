package client

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/ports/output"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/client"
)

type findClientUseCase struct {
	clientRepository repositories.ClientRepository
}

func NewFindClientUseCase(clientRepository repositories.ClientRepository) contracts.FindClientUseCase {

	return &findClientUseCase{
		clientRepository: clientRepository,
	}
}

func (c *findClientUseCase) Execute(ctx context.Context, partnerID int, name string) (*output.FindClient, error) {
	if partnerID == 0 || name == "" {
		return nil, fmt.Errorf("failed brand or name are empty")
	}

	clientEntity, err := c.clientRepository.FindClientByPartnerID(ctx, partnerID, name)
	if err != nil {
		return nil, fmt.Errorf("erro to find client with brand: '%d' at database: '%v'", partnerID, err)
	}

	if len(clientEntity) == 0 || clientEntity[0].ID == 0 {
		return nil, fmt.Errorf("client not found")
	}

	output := &output.FindClient{
		Client: clientEntity[0],
	}

	return output, nil
}
