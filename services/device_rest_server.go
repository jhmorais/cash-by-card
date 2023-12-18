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

	return router
}
