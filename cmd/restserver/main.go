package main

import (
	"fmt"
	"net/http"

	"github.com/jhmorais/cash-by-card/config"
	"github.com/jhmorais/cash-by-card/internal/infra/di"
	"github.com/jhmorais/cash-by-card/services"
)

func main() {
	config.LoadServerEnvironmentVars()

	dependencies := di.NewBuild()

	router := services.NewHTTPRouterClient(
		dependencies.Usecases.CreateClientUseCase,
		dependencies.Usecases.DeleteClientUseCase,
		dependencies.Usecases.FindClientByIDUseCase,
		dependencies.Usecases.FindClientByNameUseCase,
		dependencies.Usecases.ListClientUseCase,
		dependencies.Usecases.UpdateClientUseCase,
		dependencies.Usecases.CreatePartnerUseCase,
		dependencies.Usecases.DeletePartnerUseCase,
		dependencies.Usecases.FindPartnerByIDUseCase,
		dependencies.Usecases.FindPartnerByNameUseCase,
		dependencies.Usecases.ListPartnerUseCase,
		dependencies.Usecases.UpdatePartnerUseCase,
		dependencies.Usecases.CreateCardUseCase,
		dependencies.Usecases.DeleteCardUseCase,
		dependencies.Usecases.FindCardByIDUseCase,
		dependencies.Usecases.FindCardByLoanIDUseCase,
		dependencies.Usecases.UpdateCardUseCase,
		dependencies.Usecases.ListCardUseCase,
		dependencies.Usecases.CreateCardMachineUseCase,
		dependencies.Usecases.DeleteCardMachineUseCase,
		dependencies.Usecases.FindCardMachineByIDUseCase,
		dependencies.Usecases.UpdateCardMachineUseCase,
		dependencies.Usecases.ListCardMachineUseCase,
		dependencies.Usecases.CreateLoanUseCase,
		dependencies.Usecases.DeleteLoanUseCase,
		dependencies.Usecases.FindLoanByIDUseCase,
		dependencies.Usecases.FindLoanByClientIDUseCase,
		dependencies.Usecases.FindLoanByParnterIDUseCase,
		dependencies.Usecases.UpdateLoanUseCase,
		dependencies.Usecases.ListLoanUseCase,
		dependencies.Usecases.UpdateLoanPaymentStatusUseCase,
	)

	fmt.Println("Starting SERVER, LISTEN PORT: " + config.GetServerPort())
	clientErr := http.ListenAndServe(fmt.Sprintf(":%s", config.GetServerPort()), router)
	if clientErr != nil && clientErr != http.ErrServerClosed {
		fmt.Println("failed to create server rest on port: " + config.GetServerPort())
		fmt.Println(clientErr.Error())
	}
}
