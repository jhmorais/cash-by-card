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
	)

	fmt.Println("Starting SERVER, LISTEN PORT: " + config.GetServerPort())
	clientErr := http.ListenAndServe(fmt.Sprintf(":%s", config.GetServerPort()), router)
	if clientErr != nil && clientErr != http.ErrServerClosed {
		fmt.Println("failed to create server rest on port: " + config.GetServerPort())
		fmt.Println(clientErr.Error())
	}
}
