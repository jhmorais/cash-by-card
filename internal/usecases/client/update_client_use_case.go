package client

import (
	"context"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/client"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/client"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/client"
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
	if updateClient.Name == "" {
		return nil, fmt.Errorf("failed name client is empty")
	}

	if updateClient.PixKey == "" {
		return nil, fmt.Errorf("failed pix key client is empty")
	}

	client, err := c.clientRepository.FindClientByCPF(ctx, updateClient.CPF)
	if err != nil {
		return nil, fmt.Errorf("failed to get client")
	}

	if len(client) > 0 && client[0].ID != updateClient.ID {
		return nil, fmt.Errorf("failed, already exists client with the same cpf")
	}

	if len(updateClient.PixKey) > 250 {
		updateClient.PixKey = updateClient.PixKey[:250]
	}

	clientEntity := &entities.Client{
		ID:        updateClient.ID,
		Name:      updateClient.Name,
		PixKey:    updateClient.PixKey,
		PixType:   updateClient.PixType,
		PartnerID: updateClient.PartnerID,
		Documents: updateClient.Documents,
		Phone:     updateClient.Phone,
		CPF:       updateClient.CPF,
		UpdatedAt: time.Now(),
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
