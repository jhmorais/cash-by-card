package di

import (
	"log"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/repositories"
	"github.com/jhmorais/cash-by-card/internal/usecases/client"
	"gorm.io/gorm"
)

type DenpencyBuild struct {
	DB           *gorm.DB
	Repositories Repositories
	Usecases     Usecases
}

type Repositories struct {
	ClientRepository repositories.ClientRepository
}

type Usecases struct {
	CreateClientUseCase     contracts.CreateClientUseCase
	DeleteClientUseCase     contracts.DeleteClientUseCase
	FindClientUseCase       contracts.FindClientUseCase
	FindClientByNameUseCase contracts.FindClientByNameUseCase
	FindClientByIDUseCase   contracts.FindClientByIDUseCase
	ListClientUseCase       contracts.ListClientUseCase
	UpdateClientUseCase     contracts.UpdateClientUseCase
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
	d.Repositories.ClientRepository = repositories.NewClientRepository(d.DB)
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

	return d
}
