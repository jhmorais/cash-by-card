package client

import (
	"context"
	"fmt"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	"github.com/jhmorais/cash-by-card/internal/ports/output"
	"github.com/jhmorais/cash-by-card/internal/repositories"
)

type listClientUseCase struct {
	clientRepository repositories.ClientRepository
}

func NewListClientUseCase(clientRepository repositories.ClientRepository) contracts.ListClientUseCase {

	return &listClientUseCase{
		clientRepository: clientRepository,
	}
}

func (l *listClientUseCase) Execute(ctx context.Context) (*output.ListClient, error) {
	var err error
	output := &output.ListClient{Clients: []*entities.Client{}}

	output.Clients, err = l.clientRepository.ListClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("error when list clients on database: %v", err)
	}

	return output, nil
}
