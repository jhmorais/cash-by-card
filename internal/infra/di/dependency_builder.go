package di

import (
	"log"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	repoCard "github.com/jhmorais/cash-by-card/internal/repositories/card"
	repoClient "github.com/jhmorais/cash-by-card/internal/repositories/client"
	repoPartner "github.com/jhmorais/cash-by-card/internal/repositories/partner"
	"github.com/jhmorais/cash-by-card/internal/usecases/card"
	"github.com/jhmorais/cash-by-card/internal/usecases/client"
	"github.com/jhmorais/cash-by-card/internal/usecases/partner"
	"gorm.io/gorm"
)

type DenpencyBuild struct {
	DB           *gorm.DB
	Repositories Repositories
	Usecases     Usecases
}

type Repositories struct {
	ClientRepository  repoClient.ClientRepository
	PartnerRepository repoPartner.PartnerRepository
	CardRepository    repoCard.CardRepository
}

type Usecases struct {
	CreateClientUseCase     contracts.CreateClientUseCase
	DeleteClientUseCase     contracts.DeleteClientUseCase
	FindClientUseCase       contracts.FindClientUseCase
	FindClientByNameUseCase contracts.FindClientByNameUseCase
	FindClientByIDUseCase   contracts.FindClientByIDUseCase
	ListClientUseCase       contracts.ListClientUseCase
	UpdateClientUseCase     contracts.UpdateClientUseCase

	CreatePartnerUseCase     contracts.CreatePartnerUseCase
	DeletePartnerUseCase     contracts.DeletePartnerUseCase
	FindPartnerUseCase       contracts.FindPartnerUseCase
	FindPartnerByNameUseCase contracts.FindPartnerByNameUseCase
	FindPartnerByIDUseCase   contracts.FindPartnerByIDUseCase
	ListPartnerUseCase       contracts.ListPartnerUseCase
	UpdatePartnerUseCase     contracts.UpdatePartnerUseCase

	CreateCardUseCase       contracts.CreateCardUseCase
	DeleteCardUseCase       contracts.DeleteCardUseCase
	FindCardByIDUseCase     contracts.FindCardByIDUseCase
	FindCardByLoanIDUseCase contracts.FindCardByLoanIDUseCase
	UpdateCardUseCase       contracts.UpdateCardUseCase
	ListCardUseCase         contracts.ListCardUseCase
}

func NewBuild() *DenpencyBuild {

	builder := &DenpencyBuild{}

	builder = builder.buildDB().
		buildRepositories().
		buildUseCases()

	return builder
}

func (d *DenpencyBuild) buildDB() *DenpencyBuild {
	var err error
	d.DB, err = InitGormMysqlDB()
	if err != nil {
		log.Fatal(err)
	}
	return d
}

func (d *DenpencyBuild) buildRepositories() *DenpencyBuild {
	d.Repositories.ClientRepository = repoClient.NewClientRepository(d.DB)
	d.Repositories.PartnerRepository = repoPartner.NewPartnerRepository(d.DB)
	d.Repositories.CardRepository = repoCard.NewCardRepository(d.DB)

	return d
}

func (d *DenpencyBuild) buildUseCases() *DenpencyBuild {
	d.Usecases.CreateClientUseCase = client.NewCreateClientUseCase(d.Repositories.ClientRepository)
	d.Usecases.DeleteClientUseCase = client.NewDeleteClientUseCase(d.Repositories.ClientRepository)
	d.Usecases.FindClientUseCase = client.NewFindClientUseCase(d.Repositories.ClientRepository)
	d.Usecases.FindClientByNameUseCase = client.NewFindClientByNameUseCase(d.Repositories.ClientRepository)
	d.Usecases.FindClientByIDUseCase = client.NewFindClientByIDUseCase(d.Repositories.ClientRepository)
	d.Usecases.ListClientUseCase = client.NewListClientUseCase(d.Repositories.ClientRepository)
	d.Usecases.UpdateClientUseCase = client.NewUpdateClientUseCase(d.Repositories.ClientRepository)

	d.Usecases.CreatePartnerUseCase = partner.NewCreatePartnerUseCase(d.Repositories.PartnerRepository)
	d.Usecases.DeletePartnerUseCase = partner.NewDeletePartnerUseCase(d.Repositories.PartnerRepository)
	d.Usecases.FindPartnerUseCase = partner.NewFindPartnerUseCase(d.Repositories.PartnerRepository)
	d.Usecases.FindPartnerByNameUseCase = partner.NewFindPartnerByNameUseCase(d.Repositories.PartnerRepository)
	d.Usecases.FindPartnerByIDUseCase = partner.NewFindPartnerByIDUseCase(d.Repositories.PartnerRepository)
	d.Usecases.ListPartnerUseCase = partner.NewListPartnerUseCase(d.Repositories.PartnerRepository)
	d.Usecases.UpdatePartnerUseCase = partner.NewUpdatePartnerUseCase(d.Repositories.PartnerRepository)

	d.Usecases.CreateCardUseCase = card.NewCreateCardUseCase(d.Repositories.CardRepository)
	d.Usecases.DeleteCardUseCase = card.NewDeleteCardUseCase(d.Repositories.CardRepository)
	d.Usecases.FindCardByIDUseCase = card.NewFindCardByIDUseCase(d.Repositories.CardRepository)
	d.Usecases.FindCardByLoanIDUseCase = card.NewFindCardByLoanIDUseCase(d.Repositories.CardRepository)
	d.Usecases.ListCardUseCase = card.NewListCardsUseCase(d.Repositories.CardRepository)
	d.Usecases.UpdateCardUseCase = card.NewUpdateCardUseCase(d.Repositories.CardRepository)

	return d
}
