package client

import (
	"context"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	"github.com/jhmorais/cash-by-card/internal/ports/input"
	"github.com/jhmorais/cash-by-card/internal/ports/output"
	"github.com/jhmorais/cash-by-card/internal/repositories"
)

type updateClientUseCase struct {
	clientRepository repositories.ClientRepository
}

func NewUpdateClientUseCase(clientRepository repositories.ClientRepository) contracts.UpdateClientUseCase {

	return &updateClientUseCase{
		clientRepository: clientRepository,
	}
}

func (c *updateClientUseCase) Execute(ctx context.Context, updateClient *input.UpdateClient) (*output.CreateClient, error) {

	// if err := validator.ValidateUUId(updateClient.ID, true, "clientId"); err != nil {
	// 	return nil, err
	// }

	if updateClient.Name == "" {
		return nil, fmt.Errorf("failed name client is empty")
	}

	if updateClient.PixKey == "" {
		return nil, fmt.Errorf("failed pix key client is empty")
	}

	client, err := c.clientRepository.FindClientByPartnerID(ctx, updateClient.PartnerID, updateClient.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get client")
	}

	if len(client) > 0 && client[0].ID != 0 {
		return nil, fmt.Errorf("failed, already exists client with the same name and brand")
	}

	//max 250
	if len(updateClient.PixKey) > 250 {
		//will discard the rest
		updateClient.PixKey = updateClient.PixKey[:250]
	}

	clientEntity := &entities.Client{
		ID:        updateClient.ID,
		Name:      updateClient.Name,
		PixType:   updateClient.PixType,
		PixKey:    updateClient.PixKey,
		PartnerID: updateClient.PartnerID,
		CreatedAt: time.Now(),
	}

	errUpdate := c.clientRepository.UpdateClient(ctx, clientEntity)
	if errUpdate != nil {
		return nil, fmt.Errorf("cannot update client at database: %v", errUpdate)
	}

	createClientOutput := &output.CreateClient{
		ClientID: clientEntity.ID,
	}

	return createClientOutput, nil
}
