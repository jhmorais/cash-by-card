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
	ListPartnerUseCase      contracts.ListPartnerUseCase
}

func NewHTTPRouterClient(
	createClientUseCase contracts.CreateClientUseCase,
	deleteClientUseCase contracts.DeleteClientUseCase,
	findClientByIDUseCase contracts.FindClientByIDUseCase,
	findClientByBrandUseCase contracts.FindClientByNameUseCase,
	listClientUseCase contracts.ListClientUseCase,
	updateClientUseCase contracts.UpdateClientUseCase,
	listPartnerUseCase contracts.ListPartnerUseCase,
) *mux.Router {
	router := mux.NewRouter()
	handler := Handler{
		CreateClientUseCase:     createClientUseCase,
		DeleteClientUseCase:     deleteClientUseCase,
		FindClientByIDUseCase:   findClientByIDUseCase,
		FindClientByNameUseCase: findClientByBrandUseCase,
		ListClientUseCase:       listClientUseCase,
		UpdateClientUseCase:     updateClientUseCase,
		ListPartnerUseCase:      listPartnerUseCase,
	}
	router.UseEncodedPath()
	router.Use(utils.CommonMiddleware)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"cash-control", "content-type", "x-xsrf-token"},
	})
	router.Use(corsHandler.Handler)

	router.HandleFunc("/clients", handler.ListClients).Methods(http.MethodGet)
	router.HandleFunc("/clients/{id}", handler.GetClientByID).Methods(http.MethodGet)
	router.HandleFunc("/clients/name/{name}", handler.GetClient).Methods(http.MethodGet)
	router.HandleFunc("/clients/{id}", handler.DeleteClient).Methods(http.MethodDelete)
	router.HandleFunc("/clients", handler.CreateClient).Methods(http.MethodPost)
	router.HandleFunc("/clients/{id}", handler.UpdateClient).Methods(http.MethodPut)

	router.HandleFunc("/partners", handler.ListPartners).Methods(http.MethodGet)

	return router
}
