package di

import (
	"log"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	repoClient "github.com/jhmorais/cash-by-card/internal/repositories/client"
	repoPartner "github.com/jhmorais/cash-by-card/internal/repositories/partner"
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
}

type Usecases struct {
	CreateClientUseCase     contracts.CreateClientUseCase
	DeleteClientUseCase     contracts.DeleteClientUseCase
	FindClientUseCase       contracts.FindClientUseCase
	FindClientByNameUseCase contracts.FindClientByNameUseCase
	FindClientByIDUseCase   contracts.FindClientByIDUseCase
	ListClientUseCase       contracts.ListClientUseCase
	UpdateClientUseCase     contracts.UpdateClientUseCase

	ListPartnerUseCase contracts.ListPartnerUseCase
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

	d.Usecases.ListPartnerUseCase = partner.NewListPartnerUseCase(d.Repositories.PartnerRepository)

	return d
}
