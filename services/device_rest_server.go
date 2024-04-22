package services

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhmorais/cash-by-card/internal/contracts"
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
}

func NewHTTPRouterClient(
	createClientUseCase contracts.CreateClientUseCase,
	deleteClientUseCase contracts.DeleteClientUseCase,
	findClientByIDUseCase contracts.FindClientByIDUseCase,
	findClientByBrandUseCase contracts.FindClientByNameUseCase,
	listClientUseCase contracts.ListClientUseCase,
	updateClientUseCase contracts.UpdateClientUseCase,
	createPartnerUseCase contracts.CreatePartnerUseCase,
	deletePartnerUseCase contracts.DeletePartnerUseCase,
	findPartnerByIDUseCase contracts.FindPartnerByIDUseCase,
	findPartnerByBrandUseCase contracts.FindPartnerByNameUseCase,
	listPartnerUseCase contracts.ListPartnerUseCase,
	updatePartnerUseCase contracts.UpdatePartnerUseCase,
	createCardUseCase contracts.CreateCardUseCase,
	deleteCardUseCase contracts.DeleteCardUseCase,
	findCardByIDUseCase contracts.FindCardByIDUseCase,
	findCardByLoanIDUseCase contracts.FindCardByLoanIDUseCase,
	updateCardUseCase contracts.UpdateCardUseCase,
	listCardUseCase contracts.ListCardUseCase,
	createCardMachineUseCase contracts.CreateCardMachineUseCase,
	deleteCardMachineUseCase contracts.DeleteCardMachineUseCase,
	findCardMachineByIDUseCase contracts.FindCardMachineByIDUseCase,
	updateCardMachineUseCase contracts.UpdateCardMachineUseCase,
	listCardMachineUseCase contracts.ListCardMachineUseCase,
	createLoanUseCase contracts.CreateLoanUseCase,
	deleteLoanUseCase contracts.DeleteLoanUseCase,
	findLoanByIDUseCase contracts.FindLoanByIDUseCase,
	findLoanByClientIDUseCase contracts.FindLoanByClientIDUseCase,
	findLoanByParnterIDUseCase contracts.FindLoanByParnterIDUseCase,
	updateLoanUseCase contracts.UpdateLoanUseCase,
	listLoanUseCase contracts.ListLoanUseCase,
	updateLoanPaymentStatusUseCase contracts.UpdateLoanPaymentStatusUseCase,
) *mux.Router {
	router := mux.NewRouter()
	handler := Handler{
		CreateClientUseCase:            createClientUseCase,
		DeleteClientUseCase:            deleteClientUseCase,
		FindClientByIDUseCase:          findClientByIDUseCase,
		FindClientByNameUseCase:        findClientByBrandUseCase,
		ListClientUseCase:              listClientUseCase,
		UpdateClientUseCase:            updateClientUseCase,
		CreatePartnerUseCase:           createPartnerUseCase,
		DeletePartnerUseCase:           deletePartnerUseCase,
		FindPartnerByIDUseCase:         findPartnerByIDUseCase,
		FindPartnerByNameUseCase:       findPartnerByBrandUseCase,
		ListPartnerUseCase:             listPartnerUseCase,
		UpdatePartnerUseCase:           updatePartnerUseCase,
		CreateCardUseCase:              createCardUseCase,
		DeleteCardUseCase:              deleteCardUseCase,
		FindCardByIDUseCase:            findCardByIDUseCase,
		FindCardByLoanIDUseCase:        findCardByLoanIDUseCase,
		UpdateCardUseCase:              updateCardUseCase,
		ListCardUseCase:                listCardUseCase,
		CreateCardMachineUseCase:       createCardMachineUseCase,
		DeleteCardMachineUseCase:       deleteCardMachineUseCase,
		FindCardMachineByIDUseCase:     findCardMachineByIDUseCase,
		UpdateCardMachineUseCase:       updateCardMachineUseCase,
		ListCardMachineUseCase:         listCardMachineUseCase,
		CreateLoanUseCase:              createLoanUseCase,
		DeleteLoanUseCase:              deleteLoanUseCase,
		FindLoanByIDUseCase:            findLoanByIDUseCase,
		FindLoanByClientIDUseCase:      findLoanByClientIDUseCase,
		FindLoanByParnterIDUseCase:     findLoanByParnterIDUseCase,
		UpdateLoanUseCase:              updateLoanUseCase,
		ListLoanUseCase:                listLoanUseCase,
		UpdateLoanPaymentStatusUseCase: updateLoanPaymentStatusUseCase,
	}
	router.UseEncodedPath()
	router.Use(utils.CommonMiddleware)

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

	return router
}
