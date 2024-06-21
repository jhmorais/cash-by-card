package services

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/internal/infra/di"
	"github.com/jhmorais/cash-by-card/utils"
)

type Handler struct {
	WorkerPort              string
	CreateClientUseCase     contracts.CreateClientUseCase
	DeleteClientUseCase     contracts.DeleteClientUseCase
	FindClientByNameUseCase contracts.FindClientByNameUseCase
	FindClientByIDUseCase   contracts.FindClientByIDUseCase
	ListClientUseCase       contracts.ListClientUseCase
	UpdateClientUseCase     contracts.UpdateClientUseCase

	CreatePartnerUseCase     contracts.CreatePartnerUseCase
	DeletePartnerUseCase     contracts.DeletePartnerUseCase
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

	CreateCardMachineUseCase   contracts.CreateCardMachineUseCase
	DeleteCardMachineUseCase   contracts.DeleteCardMachineUseCase
	FindCardMachineByIDUseCase contracts.FindCardMachineByIDUseCase
	UpdateCardMachineUseCase   contracts.UpdateCardMachineUseCase
	ListCardMachineUseCase     contracts.ListCardMachineUseCase

	CreateLoanUseCase              contracts.CreateLoanUseCase
	DeleteLoanUseCase              contracts.DeleteLoanUseCase
	FindLoanByIDUseCase            contracts.FindLoanByIDUseCase
	FindLoanByClientIDUseCase      contracts.FindLoanByClientIDUseCase
	FindLoanByParnterIDUseCase     contracts.FindLoanByParnterIDUseCase
	UpdateLoanUseCase              contracts.UpdateLoanUseCase
	ListLoanUseCase                contracts.ListLoanUseCase
	UpdateLoanPaymentStatusUseCase contracts.UpdateLoanPaymentStatusUseCase

	CreateUserUseCase                 contracts.CreateUserUseCase
	LoginUseCase                      contracts.LoginUseCase
	FindUserByEmailAndPasswordUseCase contracts.FindUserByEmailAndPasswordUseCase
}

func NewHTTPRouterClient(
	useCases di.Usecases,
) *mux.Router {
	router := mux.NewRouter()
	handler := Handler{
		CreateClientUseCase:            useCases.CreateClientUseCase,
		DeleteClientUseCase:            useCases.DeleteClientUseCase,
		FindClientByIDUseCase:          useCases.FindClientByIDUseCase,
		FindClientByNameUseCase:        useCases.FindClientByNameUseCase,
		ListClientUseCase:              useCases.ListClientUseCase,
		UpdateClientUseCase:            useCases.UpdateClientUseCase,
		CreatePartnerUseCase:           useCases.CreatePartnerUseCase,
		DeletePartnerUseCase:           useCases.DeletePartnerUseCase,
		FindPartnerByIDUseCase:         useCases.FindPartnerByIDUseCase,
		FindPartnerByNameUseCase:       useCases.FindPartnerByNameUseCase,
		ListPartnerUseCase:             useCases.ListPartnerUseCase,
		UpdatePartnerUseCase:           useCases.UpdatePartnerUseCase,
		CreateCardUseCase:              useCases.CreateCardUseCase,
		DeleteCardUseCase:              useCases.DeleteCardUseCase,
		FindCardByIDUseCase:            useCases.FindCardByIDUseCase,
		FindCardByLoanIDUseCase:        useCases.FindCardByLoanIDUseCase,
		UpdateCardUseCase:              useCases.UpdateCardUseCase,
		ListCardUseCase:                useCases.ListCardUseCase,
		CreateCardMachineUseCase:       useCases.CreateCardMachineUseCase,
		DeleteCardMachineUseCase:       useCases.DeleteCardMachineUseCase,
		FindCardMachineByIDUseCase:     useCases.FindCardMachineByIDUseCase,
		UpdateCardMachineUseCase:       useCases.UpdateCardMachineUseCase,
		ListCardMachineUseCase:         useCases.ListCardMachineUseCase,
		CreateLoanUseCase:              useCases.CreateLoanUseCase,
		DeleteLoanUseCase:              useCases.DeleteLoanUseCase,
		FindLoanByIDUseCase:            useCases.FindLoanByIDUseCase,
		FindLoanByParnterIDUseCase:     useCases.FindLoanByParnterIDUseCase,
		UpdateLoanUseCase:              useCases.UpdateLoanUseCase,
		ListLoanUseCase:                useCases.ListLoanUseCase,
		UpdateLoanPaymentStatusUseCase: useCases.UpdateLoanPaymentStatusUseCase,
	}
	router.UseEncodedPath()
	router.Use(utils.CommonMiddleware)

	router.Use(utils.ValidateJwtTokenMiddleware)
	// Criar um middleware para validar o token jwt
	//		https://github.com/HunCoding/meu-primeiro-crud-go/blob/main/src/model/user_token_domain.go#L67

	router.HandleFunc("/clients", handler.ListClients).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/clients/{id}", handler.GetClientByID).Methods(http.MethodGet)
	router.HandleFunc("/clients/name/{name}", handler.GetClient).Methods(http.MethodGet)
	router.HandleFunc("/clients/{id}", handler.DeleteClient).Methods(http.MethodDelete)
	router.HandleFunc("/clients", handler.CreateClient).Methods(http.MethodPost)
	router.HandleFunc("/clients/files/{cpf}", handler.CreateClientDocuments).Methods(http.MethodPost)
	router.HandleFunc("/clients/{id}", handler.UpdateClient).Methods(http.MethodPut)

	router.HandleFunc("/partners", handler.ListPartners).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/partners/{id}", handler.GetPartnerByID).Methods(http.MethodGet)
	router.HandleFunc("/partners/name/{name}", handler.GetPartner).Methods(http.MethodGet)
	router.HandleFunc("/partners/{id}", handler.DeletePartner).Methods(http.MethodDelete)
	router.HandleFunc("/partners", handler.CreatePartner).Methods(http.MethodPost)
	router.HandleFunc("/partners/{id}", handler.UpdatePartner).Methods(http.MethodPut)

	router.HandleFunc("/cards", handler.ListCards).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/cards/{id}", handler.GetCardByID).Methods(http.MethodGet)
	router.HandleFunc("/cards/loan/{loanId}", handler.GetCardByLoanID).Methods(http.MethodGet)
	router.HandleFunc("/cards/{id}", handler.DeleteCard).Methods(http.MethodDelete)
	router.HandleFunc("/cards", handler.CreateCard).Methods(http.MethodPost)
	router.HandleFunc("/cards", handler.UpdateCard).Methods(http.MethodPut)

	router.HandleFunc("/card-machines", handler.ListCardMachines).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/card-machines/{id}", handler.GetCardMachineByID).Methods(http.MethodGet)
	router.HandleFunc("/card-machines/{id}", handler.DeleteCardMachine).Methods(http.MethodDelete)
	router.HandleFunc("/card-machines", handler.CreateCardMachine).Methods(http.MethodPost)
	router.HandleFunc("/card-machines/{id}", handler.UpdateCardMachine).Methods(http.MethodPut)

	router.HandleFunc("/loans", handler.ListLoans).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/loans/{id}", handler.GetLoanByID).Methods(http.MethodGet)
	router.HandleFunc("/loans/client/{clientId}", handler.GetLoanByClientID).Methods(http.MethodGet)
	router.HandleFunc("/loans/partner/{parnterId}", handler.GetLoanByPartnerID).Methods(http.MethodGet)
	router.HandleFunc("/loans/{id}", handler.DeleteLoan).Methods(http.MethodDelete)
	router.HandleFunc("/loans", handler.CreateLoan).Methods(http.MethodPost)
	router.HandleFunc("/loans/{id}", handler.UpdateLoan).Methods(http.MethodPut)
	router.HandleFunc("/loans/{id}/payment-status", handler.UpdateLoanPaymentStatus).Methods(http.MethodPatch)

	router.HandleFunc("/users", handler.CreateUser).Methods(http.MethodPost) // Criar o service do usuario
	router.HandleFunc("/login", handler.LoginUser).Methods(http.MethodPost)  // Criar o service do login

	return router
}
