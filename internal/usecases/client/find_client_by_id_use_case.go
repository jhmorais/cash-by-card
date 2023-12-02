package client

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/ports/output"
	"github.com/jhmorais/cash-by-card/internal/repositories"
	"github.com/jhmorais/cash-by-card/internal/usecases/validator"
)

type findClientByIDUseCase struct {
	clientRepository repositories.ClientRepository
}

func NewFindClientByIDUseCase(clientRepository repositories.ClientRepository) contracts.FindClientByIDUseCase {

	return &findClientByIDUseCase{
		clientRepository: clientRepository,
	}
}

func (c *findClientByIDUseCase) Execute(ctx context.Context, clientID string) (*output.FindClient, error) {

	if err := validator.ValidateUUId(clientID, true, "clientId"); err != nil {
		return nil, err
	}

	clientEntity, err := c.clientRepository.FindClientByID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("erro to find client '%s' at database: '%v'", clientID, err)
	}

	if clientEntity == nil || clientEntity.ID == 0 {
		return nil, fmt.Errorf("clientID not found")
	}

	output := &output.FindClient{
		Client: clientEntity,
	}

	return output, nil
}
