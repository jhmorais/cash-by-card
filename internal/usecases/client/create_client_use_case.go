package client

import (
	"context"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	"github.com/jhmorais/cash-by-card/internal/ports/input"
	"github.com/jhmorais/cash-by-card/internal/ports/output"
	repositories "github.com/jhmorais/cash-by-card/internal/repositories/client"
)

type createClientUseCase struct {
	clientRepository repositories.ClientRepository
}

func NewCreateClientUseCase(clientRepository repositories.ClientRepository) contracts.CreateClientUseCase {

	return &createClientUseCase{
		clientRepository: clientRepository,
	}
}

func (c *createClientUseCase) Execute(ctx context.Context, createClient *input.CreateClient) (*output.CreateClient, error) {

	if len(createClient.Name) > 250 {
		createClient.Name = createClient.Name[:250]
	}

	if createClient.Phone == "" {
		return nil, fmt.Errorf("cannot create a client without phone")
	}

	if createClient.CPF == "" {
		return nil, fmt.Errorf("cannot create a client without cpf")
	}

	if len(createClient.Documents) > 100 {
		return nil, fmt.Errorf("cannot have documents greater than 100 characters")
	}

	client, err := c.clientRepository.FindClientByName(ctx, createClient.Name)
	if err != nil {
		return nil, fmt.Errorf("failed to get client")
	}

	if len(client) > 0 && client[0].ID > 0 {
		return nil, fmt.Errorf("failed, already exists client with the same name")
	}

	client, err = c.clientRepository.FindClientByCPF(ctx, createClient.CPF)
	if err != nil {
		return nil, fmt.Errorf("failed to get client: %v", err)
	}

	if len(client) > 0 {
		return nil, fmt.Errorf("failed, already exists client with the same cpf")
	}

	clientEntity := &entities.Client{
		Name:      createClient.Name,
		PixType:   createClient.PixType,
		PixKey:    createClient.PixKey,
		PartnerID: createClient.PartnerID,
		CPF:       createClient.CPF,
		Phone:     createClient.Phone,
		Documents: createClient.Documents,
		CreatedAt: time.Now(),
	}

	err = c.clientRepository.CreateClient(ctx, clientEntity)
	if err != nil {
		return nil, fmt.Errorf("cannot save client at database: %v", err)
	}

	createClientOutput := &output.CreateClient{
		ClientID: clientEntity.ID,
	}

	return createClientOutput, nil
}
