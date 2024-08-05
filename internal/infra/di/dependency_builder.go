package di

import (
	"log"

	"github.com/jhmorais/cash-by-card/internal/contracts"
	repoCard "github.com/jhmorais/cash-by-card/internal/repositories/card"
	repoCardMachine "github.com/jhmorais/cash-by-card/internal/repositories/cardMachine"
	repoClient "github.com/jhmorais/cash-by-card/internal/repositories/client"
	repoLoan "github.com/jhmorais/cash-by-card/internal/repositories/loan"
	repoPartner "github.com/jhmorais/cash-by-card/internal/repositories/partner"
	repoUser "github.com/jhmorais/cash-by-card/internal/repositories/user"
	"github.com/jhmorais/cash-by-card/internal/usecases/card"
	"github.com/jhmorais/cash-by-card/internal/usecases/cardMachine"
	"github.com/jhmorais/cash-by-card/internal/usecases/client"
	"github.com/jhmorais/cash-by-card/internal/usecases/loan"
	"github.com/jhmorais/cash-by-card/internal/usecases/login"
	"github.com/jhmorais/cash-by-card/internal/usecases/partner"
	"github.com/jhmorais/cash-by-card/internal/usecases/user"
	"gorm.io/gorm"
)

type DenpencyBuild struct {
	DB           *gorm.DB
	Repositories Repositories
	Usecases     Usecases
}

type Repositories struct {
	ClientRepository      repoClient.ClientRepository
	PartnerRepository     repoPartner.PartnerRepository
	CardRepository        repoCard.CardRepository
	CardMachineRepository repoCardMachine.CardMachineRepository
	LoanRepository        repoLoan.LoanRepository
	UserRepository        repoUser.UserRepository
}

type Usecases struct {
	CreateClientUseCase     contracts.CreateClientUseCase
	DeleteClientUseCase     contracts.DeleteClientUseCase
	FindClientUseCase       contracts.FindClientUseCase
	FindClientByNameUseCase contracts.FindClientByNameUseCase
	FindClientByIDUseCase   contracts.FindClientByIDUseCase
	ListClientUseCase       contracts.ListClientUseCase
	UpdateClientUseCase     contracts.UpdateClientUseCase

	CreatePartnerUseCase      contracts.CreatePartnerUseCase
	DeletePartnerUseCase      contracts.DeletePartnerUseCase
	FindPartnerUseCase        contracts.FindPartnerUseCase
	FindPartnerByNameUseCase  contracts.FindPartnerByNameUseCase
	FindPartnerByEmailUseCase contracts.FindPartnerByEmailUseCase
	FindPartnerByIDUseCase    contracts.FindPartnerByIDUseCase
	ListPartnerUseCase        contracts.ListPartnerUseCase
	UpdatePartnerUseCase      contracts.UpdatePartnerUseCase

	CreateCardUseCase       contracts.CreateCardUseCase
	DeleteCardUseCase       contracts.DeleteCardUseCase
	FindCardByIDUseCase     contracts.FindCardByIDUseCase
	FindCardByLoanIDUseCase contracts.FindCardByLoanIDUseCase
	UpdateCardUseCase       contracts.UpdateCardUseCase
	ListCardUseCase         contracts.ListCardUseCase

	CreateCardMachineUseCase   contracts.CreateCardMachineUseCase
	DeleteCardMachineUseCase   contracts.DeleteCardMachineUseCase
	FindCardMachineByIDUseCase contracts.FindCardMachineByIDUseCase
	UpdateCardMachineUseCase   contracts.UpdateCardMachineUseCase
	ListCardMachineUseCase     contracts.ListCardMachineUseCase

	CreateLoanUseCase              contracts.CreateLoanUseCase
	DeleteLoanUseCase              contracts.DeleteLoanUseCase
	FindLoanByIDUseCase            contracts.FindLoanByIDUseCase
	UpdateLoanUseCase              contracts.UpdateLoanUseCase
	ListLoanUseCase                contracts.ListLoanUseCase
	FindLoanByClientIDUseCase      contracts.FindLoanByClientIDUseCase
	FindLoanByParnterIDUseCase     contracts.FindLoanByParnterIDUseCase
	UpdateLoanPaymentStatusUseCase contracts.UpdateLoanPaymentStatusUseCase

	CreateUserUseCase                 contracts.CreateUserUseCase
	LoginUseCase                      contracts.LoginUseCase
	FindUserByEmailAndPasswordUseCase contracts.FindUserByEmailAndPasswordUseCase
	ListUserUseCase                   contracts.ListUserUseCase
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
	d.Repositories.CardMachineRepository = repoCardMachine.NewCardMachineRepository(d.DB)
	d.Repositories.LoanRepository = repoLoan.NewLoanRepository(d.DB)
	d.Repositories.UserRepository = repoUser.NewUserRepository(d.DB)

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
	d.Usecases.FindPartnerByEmailUseCase = partner.NewFindPartnerByEmailUseCase(d.Repositories.PartnerRepository)
	d.Usecases.FindPartnerByIDUseCase = partner.NewFindPartnerByIDUseCase(d.Repositories.PartnerRepository)
	d.Usecases.ListPartnerUseCase = partner.NewListPartnerUseCase(d.Repositories.PartnerRepository)
	d.Usecases.UpdatePartnerUseCase = partner.NewUpdatePartnerUseCase(d.Repositories.PartnerRepository)

	d.Usecases.CreateCardUseCase = card.NewCreateCardUseCase(d.Repositories.CardRepository)
	d.Usecases.DeleteCardUseCase = card.NewDeleteCardUseCase(d.Repositories.CardRepository)
	d.Usecases.FindCardByIDUseCase = card.NewFindCardByIDUseCase(d.Repositories.CardRepository)
	d.Usecases.FindCardByLoanIDUseCase = card.NewFindCardByLoanIDUseCase(d.Repositories.CardRepository)
	d.Usecases.ListCardUseCase = card.NewListCardsUseCase(d.Repositories.CardRepository)
	d.Usecases.UpdateCardUseCase = card.NewUpdateCardUseCase(d.Repositories.CardRepository)

	d.Usecases.CreateCardMachineUseCase = cardMachine.NewCreateCardMachineUseCase(d.Repositories.CardMachineRepository)
	d.Usecases.DeleteCardMachineUseCase = cardMachine.NewDeleteCardMachineUseCase(d.Repositories.CardMachineRepository)
	d.Usecases.FindCardMachineByIDUseCase = cardMachine.NewFindCardMachineByIDUseCase(d.Repositories.CardMachineRepository)
	d.Usecases.ListCardMachineUseCase = cardMachine.NewListCardMachinesUseCase(d.Repositories.CardMachineRepository)
	d.Usecases.UpdateCardMachineUseCase = cardMachine.NewUpdateCardMachineUseCase(d.Repositories.CardMachineRepository)

	d.Usecases.DeleteLoanUseCase = loan.NewDeleteLoanUseCase(d.Repositories.LoanRepository)
	d.Usecases.FindLoanByIDUseCase = loan.NewFindLoanByIDUseCase(d.Repositories.LoanRepository)
	d.Usecases.FindLoanByClientIDUseCase = loan.NewFindLoanByClientIDUseCase(d.Repositories.LoanRepository)
	d.Usecases.FindLoanByParnterIDUseCase = loan.NewFindLoansByPartnerIDUseCase(d.Repositories.LoanRepository)
	d.Usecases.ListLoanUseCase = loan.NewListLoansUseCase(d.Repositories.LoanRepository)
	d.Usecases.UpdateLoanUseCase = loan.NewUpdateLoanUseCase(d.Repositories.LoanRepository)
	d.Usecases.CreateLoanUseCase = loan.NewCreateLoanUseCase(d.Repositories.LoanRepository, d.Usecases.CreateCardUseCase)
	d.Usecases.UpdateLoanPaymentStatusUseCase = loan.NewUpdateLoanPaymentStatusUseCase(d.Repositories.LoanRepository)

	d.Usecases.CreateUserUseCase = user.NewCreateUserUseCase(d.Repositories.UserRepository)
	d.Usecases.FindUserByEmailAndPasswordUseCase = user.NewFindUserByEmailAndPasswordUseCase(d.Repositories.UserRepository)
	d.Usecases.LoginUseCase = login.NewLoginUseCase(d.Repositories.UserRepository)

	return d
}
