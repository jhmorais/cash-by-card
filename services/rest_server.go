package services

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhmorais/cash-by-card/internal/infra/di"
	"github.com/jhmorais/cash-by-card/utils"
)

type Handler struct {
	WorkerPort string
	Usecases   di.Usecases
}

func NewHTTPRouterClient(
	useCases di.Usecases,
	userRepo di.Repositories,
) *mux.Router {
	router := mux.NewRouter()
	handler := Handler{Usecases: useCases}
	router.UseEncodedPath()
	router.Use(utils.CommonMiddleware)

	adminRouter := router.PathPrefix("/admin").Subrouter()
	publicRouter := router.PathPrefix("/public").Subrouter()
	authRouter := router.PathPrefix("/auth").Subrouter()

	publicRouter.Use(utils.ValidateJwtTokenMiddleware)
	adminRouter.Use(utils.ValidateJwtTokenMiddleware)
	adminRouter.Use(utils.RoleMiddleware("admin", userRepo.UserRepository))

	adminRouter.HandleFunc("/clients", handler.ListClients).Methods(http.MethodGet, http.MethodOptions)
	adminRouter.HandleFunc("/clients/{id}", handler.GetClientByID).Methods(http.MethodGet)
	adminRouter.HandleFunc("/clients/name/{name}", handler.GetClient).Methods(http.MethodGet)
	adminRouter.HandleFunc("/clients/{id}", handler.DeleteClient).Methods(http.MethodDelete)
	adminRouter.HandleFunc("/clients", handler.CreateClient).Methods(http.MethodPost)
	adminRouter.HandleFunc("/clients/files/{cpf}", handler.CreateClientDocuments).Methods(http.MethodPost)
	adminRouter.HandleFunc("/clients/{id}", handler.UpdateClient).Methods(http.MethodPut)

	publicRouter.HandleFunc("/partners", handler.ListPartners).Methods(http.MethodGet, http.MethodOptions)
	publicRouter.HandleFunc("/partners/{id}", handler.GetPartnerByID).Methods(http.MethodGet)
	publicRouter.HandleFunc("/partners/name/{name}", handler.GetPartner).Methods(http.MethodGet)
	publicRouter.HandleFunc("/partners/{id}", handler.DeletePartner).Methods(http.MethodDelete)
	publicRouter.HandleFunc("/partners", handler.CreatePartner).Methods(http.MethodPost)
	publicRouter.HandleFunc("/partners/{id}", handler.UpdatePartner).Methods(http.MethodPut)

	adminRouter.HandleFunc("/cards", handler.ListCards).Methods(http.MethodGet, http.MethodOptions)
	adminRouter.HandleFunc("/cards/{id}", handler.GetCardByID).Methods(http.MethodGet)
	adminRouter.HandleFunc("/cards/loan/{loanId}", handler.GetCardByLoanID).Methods(http.MethodGet)
	adminRouter.HandleFunc("/cards/{id}", handler.DeleteCard).Methods(http.MethodDelete)
	adminRouter.HandleFunc("/cards", handler.CreateCard).Methods(http.MethodPost)
	adminRouter.HandleFunc("/cards", handler.UpdateCard).Methods(http.MethodPut)

	adminRouter.HandleFunc("/card-machines", handler.ListCardMachines).Methods(http.MethodGet, http.MethodOptions)
	adminRouter.HandleFunc("/card-machines/{id}", handler.GetCardMachineByID).Methods(http.MethodGet)
	adminRouter.HandleFunc("/card-machines/{id}", handler.DeleteCardMachine).Methods(http.MethodDelete)
	adminRouter.HandleFunc("/card-machines", handler.CreateCardMachine).Methods(http.MethodPost)
	adminRouter.HandleFunc("/card-machines/{id}", handler.UpdateCardMachine).Methods(http.MethodPut)

	adminRouter.HandleFunc("/loans", handler.ListLoans).Methods(http.MethodGet, http.MethodOptions)
	adminRouter.HandleFunc("/loans/{id}", handler.GetLoanByID).Methods(http.MethodGet)
	adminRouter.HandleFunc("/loans/client/{clientId}", handler.GetLoanByClientID).Methods(http.MethodGet)
	adminRouter.HandleFunc("/loans/partner/{parnterId}", handler.GetLoanByPartnerID).Methods(http.MethodGet)
	adminRouter.HandleFunc("/loans/{id}", handler.DeleteLoan).Methods(http.MethodDelete)
	adminRouter.HandleFunc("/loans", handler.CreateLoan).Methods(http.MethodPost)
	adminRouter.HandleFunc("/loans/{id}", handler.UpdateLoan).Methods(http.MethodPut)
	adminRouter.HandleFunc("/loans/{id}/payment-status", handler.UpdateLoanPaymentStatus).Methods(http.MethodPatch)

	adminRouter.HandleFunc("/users", handler.CreateUser).Methods(http.MethodPost)
	adminRouter.HandleFunc("/users", handler.ListUsers).Methods(http.MethodGet)
	adminRouter.HandleFunc("/users", handler.UpdateUser).Methods(http.MethodPut)
	adminRouter.HandleFunc("/users", handler.SendCode).Methods(http.MethodPost)
	authRouter.HandleFunc("/login", handler.LoginUser).Methods(http.MethodPost)

	return router
}
