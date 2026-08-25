package partner

import (
	"context"
	"fmt"
	"time"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/domain/entities"
	input "github.com/jhmorais/cash-by-card/internal/ports/input/client"
	output "github.com/jhmorais/cash-by-card/internal/ports/output/partner"
	repoClient "github.com/jhmorais/cash-by-card/internal/repositories/client"
	repoPartner "github.com/jhmorais/cash-by-card/internal/repositories/partner"
)

const partnerEditWindow = 24 * time.Hour

type partnerClientsUseCase struct {
	clientRepository  repoClient.ClientRepository
	partnerRepository repoPartner.PartnerRepository
}

func NewPartnerClientsUseCase(clientRepository repoClient.ClientRepository, partnerRepository repoPartner.PartnerRepository) contracts.PartnerClientsUseCase {
	return &partnerClientsUseCase{clientRepository: clientRepository, partnerRepository: partnerRepository}
}

// resolvePartner devolve a entidade parceira pelo email do user; nil, nil quando não existe.
func (p *partnerClientsUseCase) resolvePartner(ctx context.Context, email string) (*entities.Partner, error) {
	partners, err := p.partnerRepository.FindPartnerByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if len(partners) == 0 || partners[0].Email == "" {
		return nil, nil
	}
	return partners[0], nil
}

func (p *partnerClientsUseCase) ListClients(ctx context.Context, partnerUserEmail string) (*output.PartnerClients, error) {
	partner, err := p.resolvePartner(ctx, partnerUserEmail)
	if err != nil {
		return nil, err
	}
	if partner == nil {
		return &output.PartnerClients{Clients: []*entities.Client{}}, nil
	}
	clients, err := p.clientRepository.ListClientsByPartnerID(ctx, partner.ID)
	if err != nil {
		return nil, err
	}
	if clients == nil {
		clients = []*entities.Client{}
	}
	return &output.PartnerClients{Clients: clients}, nil
}

func (p *partnerClientsUseCase) CreateClient(ctx context.Context, partnerUserEmail string, createClient *input.CreateClient) (*output.PartnerCreateClient, error) {
	partner, err := p.resolvePartner(ctx, partnerUserEmail)
	if err != nil {
		return nil, err
	}
	if partner == nil {
		return nil, fmt.Errorf("nenhum parceiro vinculado a este usuário")
	}
	// vínculo automático: o cliente nasce ligado à entidade do parceiro logado
	partnerID := partner.ID
	entity := &entities.Client{
		Name:      createClient.Name,
		PixType:   createClient.PixType,
		PixKey:    createClient.PixKey,
		Phone:     createClient.Phone,
		CPF:       createClient.CPF,
		Documents: createClient.Documents,
		PartnerID: &partnerID,
		CreatedAt: time.Now(),
	}
	if err := p.clientRepository.CreateClient(ctx, entity); err != nil {
		return nil, err
	}
	return &output.PartnerCreateClient{ClientID: entity.ID}, nil
}

func (p *partnerClientsUseCase) UpdateClient(ctx context.Context, partnerUserEmail string, updateClient *input.UpdateClient) (*output.PartnerUpdateClient, error) {
	partner, err := p.resolvePartner(ctx, partnerUserEmail)
	if err != nil {
		return nil, err
	}
	if partner == nil {
		return nil, fmt.Errorf("nenhum parceiro vinculado a este usuário")
	}
	client, err := p.clientRepository.FindClientByID(ctx, updateClient.ID)
	if err != nil || client == nil || client.ID == 0 {
		return nil, fmt.Errorf("cliente não encontrado")
	}
	if client.PartnerID == nil || *client.PartnerID != partner.ID {
		return nil, fmt.Errorf("sem permissão para editar este cliente")
	}
	if time.Since(client.CreatedAt) > partnerEditWindow {
		return nil, fmt.Errorf("clientes só podem ser editados nas primeiras 24h após a criação; procure um administrador")
	}
	client.Name = updateClient.Name
	client.PixType = updateClient.PixType
	client.PixKey = updateClient.PixKey
	client.Phone = updateClient.Phone
	client.CPF = updateClient.CPF
	client.Documents = updateClient.Documents
	if err := p.clientRepository.UpdateClient(ctx, client); err != nil {
		return nil, err
	}
	return &output.PartnerUpdateClient{ClientID: client.ID}, nil
}
