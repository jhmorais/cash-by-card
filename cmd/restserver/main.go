package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jhmorais/cash-by-card/config"
	"github.com/jhmorais/cash-by-card/internal/infra/di"
	"github.com/jhmorais/cash-by-card/services"
)

type MyServer struct {
	r *mux.Router
}

func main() {
	config.LoadServerEnvironmentVars()

	dependencies := di.NewBuild()

	router := services.NewHTTPRouterClient(
		dependencies.Usecases,
	)

	http.Handle("/", &MyServer{router})

	fmt.Println("Starting SERVER, LISTEN PORT: " + config.GetServerPort())
	clientErr := http.ListenAndServe(fmt.Sprintf(":%s", config.GetServerPort()), nil)
	if clientErr != nil && clientErr != http.ErrServerClosed {
		fmt.Println("failed to create server rest on port: " + config.GetServerPort())
		fmt.Println(clientErr.Error())
	}
}

func (s *MyServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if origin := req.Header.Get("Origin"); origin != "" {
		rw.Header().Set("Access-Control-Allow-Origin", origin)
		rw.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, PATCH, DELETE")
		rw.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	}
	if req.Method == "OPTIONS" {
		return
	}
	s.r.ServeHTTP(rw, req)
}
