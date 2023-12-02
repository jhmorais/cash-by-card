package client

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/ports/output"
	"github.com/jhmorais/cash-by-card/internal/repositories"
	"github.com/jhmorais/cash-by-card/internal/usecases/validator"
)

type deleteClientUseCase struct {
	clientRepository repositories.ClientRepository
}

func NewDeleteClientUseCase(clientRepository repositories.ClientRepository) contracts.DeleteClientUseCase {

	return &deleteClientUseCase{
		clientRepository: clientRepository,
	}
}

func (c *deleteClientUseCase) Execute(ctx context.Context, clientID string) (*output.DeleteClient, error) {

	if err := validator.ValidateUUId(clientID, true, "clientId"); err != nil {
		return nil, err
	}

	clientEntity, err := c.clientRepository.FindClientByID(ctx, clientID)
	if err != nil {
		return nil, fmt.Errorf("failed to find client '%s' at database: '%v'", clientID, err)
	}

	if clientEntity == nil || clientEntity.ID == 0 {
		return nil, fmt.Errorf("clientID not found")
	}

	err = c.clientRepository.DeleteClient(ctx, clientEntity)
	if err != nil {
		return nil, fmt.Errorf("failed to delete client '%s'", clientEntity.ID)
	}

	output := &output.DeleteClient{
		ClientID:   clientEntity.ID,
		ClientName: clientEntity.Name,
	}

	return output, nil
}
