package services

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhmorais/cash-by-card/internal/contracts"
	"github.com/jhmorais/cash-by-card/utils"
	"github.com/rs/cors"
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
) *mux.Router {
	router := mux.NewRouter()
	handler := Handler{
		CreateClientUseCase:      createClientUseCase,
		DeleteClientUseCase:      deleteClientUseCase,
		FindClientByIDUseCase:    findClientByIDUseCase,
		FindClientByNameUseCase:  findClientByBrandUseCase,
		ListClientUseCase:        listClientUseCase,
		UpdateClientUseCase:      updateClientUseCase,
		CreatePartnerUseCase:     createPartnerUseCase,
		DeletePartnerUseCase:     deletePartnerUseCase,
		FindPartnerByIDUseCase:   findPartnerByIDUseCase,
		FindPartnerByNameUseCase: findPartnerByBrandUseCase,
		ListPartnerUseCase:       listPartnerUseCase,
		UpdatePartnerUseCase:     updatePartnerUseCase,
		CreateCardUseCase:        createCardUseCase,
		DeleteCardUseCase:        deleteCardUseCase,
		FindCardByIDUseCase:      findCardByIDUseCase,
		FindCardByLoanIDUseCase:  findCardByLoanIDUseCase,
		UpdateCardUseCase:        updateCardUseCase,
		ListCardUseCase:          listCardUseCase,
	}
	router.UseEncodedPath()
	router.Use(utils.CommonMiddleware)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"cash-control", "content-type", "x-xsrf-token"},
	})
	router.Use(corsHandler.Handler)

	router.HandleFunc("/clients", handler.ListClients).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/clients/{id}", handler.GetClientByID).Methods(http.MethodGet)
	router.HandleFunc("/clients/name/{name}", handler.GetClient).Methods(http.MethodGet)
	router.HandleFunc("/clients/{id}", handler.DeleteClient).Methods(http.MethodDelete)
	router.HandleFunc("/clients", handler.CreateClient).Methods(http.MethodPost)
	router.HandleFunc("/clients/{id}", handler.UpdateClient).Methods(http.MethodPut)

	router.HandleFunc("/partners", handler.ListPartners).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/partners/{id}", handler.GetPartnerByID).Methods(http.MethodGet)
	router.HandleFunc("/partners/name/{name}", handler.GetPartner).Methods(http.MethodGet)
	router.HandleFunc("/partners/{id}", handler.DeletePartner).Methods(http.MethodDelete)
	router.HandleFunc("/partners", handler.CreatePartner).Methods(http.MethodPost)
	router.HandleFunc("/partners/{id}", handler.UpdatePartner).Methods(http.MethodPut)

	router.HandleFunc("/cards", handler.ListCards).Methods(http.MethodGet, http.MethodOptions)
	router.HandleFunc("/cards/{id}", handler.GetCardByID).Methods(http.MethodGet)
	router.HandleFunc("/cards/{loanId}", handler.GetCardByLoanID).Methods(http.MethodGet)
	router.HandleFunc("/cards/{id}", handler.DeleteCard).Methods(http.MethodDelete)
	router.HandleFunc("/cards", handler.CreateCard).Methods(http.MethodPost)
	router.HandleFunc("/cards/{id}", handler.UpdateCard).Methods(http.MethodPut)

	return router
}
